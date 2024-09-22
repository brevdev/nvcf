package collections

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Filter[T any](list []T, f func(T) bool) []T {
	result := []T{}
	for _, item := range list {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func ListContains[K comparable](list []K, item K) bool {
	return ListHas(list, func(l K) bool { return l == item })
}

func ListHas[K any](list []K, has func(l K) bool) bool {
	k := Find(list, has)
	if k != nil {
		return true
	}
	return false
}

func Find[T any](list []T, f func(T) bool) *T {
	for _, item := range list {
		if f(item) {
			return &item
		}
	}
	return nil
}

func MapE[T, R any](items []T, mapper func(T) (R, error)) ([]R, error) {
	results := []R{}
	for _, item := range items {
		res, err := mapper(item)
		if err != nil {
			return results, err
		}
		results = append(results, res)
	}
	return results, nil
}

func TryCopyToNew[T any, R any](t T, options ...CopyOption) (R, error) {
	var r R
	copyOptions := CopyOptions{}
	for _, o := range options {
		o.apply(t, r, &copyOptions)
	}
	if err := copier.CopyWithOption(&r, t, copyOptions.ToCopierOptions()); err != nil {
		return r, err
	}
	return r, nil
}

type CopyOption interface {
	apply(t any, r any, o *CopyOptions)
}

type CopyMap map[string]string

func (c CopyMap) apply(t any, r any, o *CopyOptions) {
	o.Mappers = append(o.Mappers, DumbCopyMapper{Mapping: copier.FieldNameMapping{
		SrcType: t,
		DstType: r,
		Mapping: c,
	}})
}

func TryCopyToNewOptions[T any, R any](t T, options CopyOptions) (R, error) {
	var r R
	if err := copier.CopyWithOption(&r, t, options.ToCopierOptions()); err != nil {
		return r, err
	}
	return r, nil
}

// For testing different DX
func TryCopyToNewE[T any, R any](t T, mappers ...CopyMappingFunc[T, R]) (R, error) {
	var r R
	opts := CopyOptions{}.ToCopierOptions()
	opts.FieldNameMapping = append(opts.FieldNameMapping, Map(mappers, func(m CopyMappingFunc[T, R]) copier.FieldNameMapping { return m.ToCopierMapping() })...)
	if err := copier.CopyWithOption(&r, t, opts); err != nil {
		return r, err
	}
	return r, nil
}

func Map[T, R any](items []T, mapper func(T) R) []R {
	results := []R{}
	for _, item := range items {
		results = append(results, mapper(item))
	}
	return results
}

type CopyOptions struct {
	ShallowCopy           bool
	OmitDefaultConverters bool
	Converters            CopierConverters
	Mappers               CopyMappers // create mappings for any arbitrary type
}

type CopyMappers []CopyMapper

func (c CopyMappers) ToCopierMappings() []copier.FieldNameMapping {
	var result []copier.FieldNameMapping
	for _, cc := range c {
		result = append(result, cc.ToCopierMapping())
	}
	return result
}

type CopyMapper interface {
	ToCopierMapping() copier.FieldNameMapping
}

type DumbCopyMapper struct {
	Mapping copier.FieldNameMapping
}

func (d DumbCopyMapper) ToCopierMapping() copier.FieldNameMapping {
	return d.Mapping
}

type CopyMapping[T, R any] map[string]string

func (c CopyMapping[T, R]) ToCopierMapping() copier.FieldNameMapping {
	var t T
	var r R
	return copier.FieldNameMapping{
		SrcType: t,
		DstType: r,
		Mapping: c,
	}
}

type CopyMappingFunc[T, R any] func(T, R) map[any]any

// Doesn't work but leaving for future reference
func (c CopyMappingFunc[T, R]) ToCopierMapping() copier.FieldNameMapping {
	var t T
	valueT := reflect.ValueOf(&t).Elem()
	var r R
	valueR := reflect.ValueOf(&r).Elem()
	copierMapping := map[string]string{}
	realMapping := c(t, r)
	var allErr error
	for tf, rf := range realMapping {
		tfFound := findStructField(valueT, reflect.ValueOf(tf))
		rfFound := findStructField(valueR, reflect.ValueOf(rf))
		isErr := false
		if tfFound == nil {
			allErr = fmt.Errorf("%v: field %v not found in struct %T", allErr, tf, t)
			isErr = true
		}
		if rfFound == nil {
			allErr = fmt.Errorf("%v: field %v not found in struct %T", allErr, rf, r)
			isErr = true
		}
		if !isErr {
			copierMapping[tfFound.Name] = rfFound.Name
		}
	}
	if allErr != nil {
		panic(allErr)
	}
	return copier.FieldNameMapping{
		SrcType: t,
		DstType: r,
		Mapping: copierMapping,
	}
}

// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type == fieldValue.Elem().Type() {
				return &sf
			}
		}
		if sf.Anonymous {
			// delve into anonymous struct to look for the field
			fi := structValue.Field(i)
			if sf.Type.Kind() == reflect.Ptr {
				fi = fi.Elem()
			}
			if fi.Kind() == reflect.Struct {
				if f := findStructField(fi, fieldValue); f != nil {
					return f
				}
			}
		}
	}
	return nil
}

var timeToPBTimeStamp CopyConverter[time.Time, *timestamppb.Timestamp] = func(src time.Time) (*timestamppb.Timestamp, error) {
	return timestamppb.New(src), nil
}

var timePtrToPBTimeStamp CopyConverter[*time.Time, *timestamppb.Timestamp] = func(src *time.Time) (*timestamppb.Timestamp, error) {
	if src == nil {
		return nil, nil
	}
	return timestamppb.New(*src), nil
}

var pbTimeStampToTime CopyConverter[*timestamppb.Timestamp, time.Time] = func(src *timestamppb.Timestamp) (time.Time, error) {
	return src.AsTime(), nil
}

var pbTimeStampToTimePtr CopyConverter[*timestamppb.Timestamp, *time.Time] = func(src *timestamppb.Timestamp) (*time.Time, error) {
	if src == nil {
		return nil, nil
	}
	t := src.AsTime()
	return &t, nil
}

var DefaultConverters = CopierConverters{
	timeToPBTimeStamp,
	timePtrToPBTimeStamp,
	pbTimeStampToTime,
	pbTimeStampToTimePtr,
}

func (c CopyOptions) ToCopierOptions() copier.Option {
	convs := c.Converters.ToCopierTypeConverters()
	if !c.OmitDefaultConverters {
		convs = append(convs, DefaultConverters.ToCopierTypeConverters()...)
	}
	return copier.Option{
		DeepCopy:         !c.ShallowCopy,
		Converters:       convs,
		FieldNameMapping: c.Mappers.ToCopierMappings(),
	}
}

type CopierConverters []CopierConverter

func (c CopierConverters) ToCopierTypeConverters() []copier.TypeConverter {
	var result []copier.TypeConverter
	for _, cc := range c {
		result = append(result, cc.ToCopierTypeConverter())
	}
	return result
}

type CopierConverter interface {
	ToCopierTypeConverter() copier.TypeConverter
}

type CopyConverter[T, R any] func(T) (R, error)

func (c CopyConverter[T, R]) ToCopierTypeConverter() copier.TypeConverter {
	var t T
	var r R
	ctc := copier.TypeConverter{
		SrcType: t,
		DstType: r,
		Fn: func(src interface{}) (dst interface{}, err error) {
			tt, ok := src.(T)
			if !ok {
				return nil, fmt.Errorf("cannot convert %T to %T", src, dst)
			}
			return c(tt)
		},
	}
	return ctc
}

func TryCopyTo[T any, R any](t T, r R) (R, error) {
	if err := copier.CopyWithOption(&r, t, CopyOptions{}.ToCopierOptions()); err != nil {
		return r, fmt.Errorf("cannot convert %T to %T", t, r)
	}
	return r, nil
}

func TryCopyToOptions[T any, R any](t T, r R, options CopyOptions) (R, error) {
	if err := copier.CopyWithOption(&r, t, options.ToCopierOptions()); err != nil {
		return r, fmt.Errorf("cannot convert %T to %T", t, r)
	}
	return r, nil
}
