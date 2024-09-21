package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func Error(cmd *cobra.Command, message string, err error) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	var formattedError error
	if verbose {
		formattedError = fmt.Errorf("%s: %v", message, err)
	} else {
		formattedError = fmt.Errorf("%s", message)
	}
	return formattedError
}

// func Success(cmd *cobra.Command, message string) {
// 	if !isQuiet(cmd) {
// 		color.Green(message)
// 	}
// }

// func Info(cmd *cobra.Command, message string) {
// 	if !isQuiet(cmd) {
// 		color.Blue(message)
// 	}
// }

func isJSON(cmd *cobra.Command) bool {
	json, _ := cmd.Flags().GetBool("json")
	return json
}

func isQuiet(cmd *cobra.Command) bool {
	quiet, _ := cmd.Flags().GetBool("quiet")
	return quiet
}

func printJSON(cmd *cobra.Command, data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		Error(cmd, "Error formatting JSON", err)
		return
	}
	fmt.Println(string(json))
}

func Prompt(message string, isSecret bool) string {
	Type(message)
	if isSecret {
		// Implement secure input for secrets
	}
	var input string
	fmt.Scanln(&input)
	return input
}

// NewSpinner creates and returns a new spinner
func NewSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
	s.Suffix = suffix
	return s
}

// StartSpinner starts the spinner
func StartSpinner(s *spinner.Spinner) {
	s.Start()
}

// StopSpinner stops the spinner
func StopSpinner(s *spinner.Spinner) {
	s.Stop()
}

// Implement other output functions (Function, Deployment, InvocationResult, etc.) here
func Functions(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, functions)
	} else {
		printFunctionsTable(cmd, functions)
	}
}

func printFunctionsTable(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	for _, f := range functions {
		table.Append([]string{f.Name, f.VersionID, string(f.Status)})
	}
	table.Render()
}

func SingleFunction(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, fn)
	} else {
		printSingleFunctionTable(cmd, fn)
	}
}

func printSingleFunctionTable(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{fn.Name, fn.VersionID, string(fn.Status)})
	table.Render()
}

func Deployments(cmd *cobra.Command, deployments []nvcf.DeploymentResponse) {
	if isJSON(cmd) {
		printJSON(cmd, deployments)
	} else {
		printDeploymentsTable(cmd, deployments)
	}
}

func printDeploymentsTable(cmd *cobra.Command, deployments []nvcf.DeploymentResponse) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Function ID", "Function Version ID", "Status"})
	table.SetBorder(false)
	for _, deployment := range deployments {
		table.Append([]string{deployment.Deployment.FunctionID, deployment.Deployment.FunctionVersionID, string(deployment.Deployment.FunctionStatus)})
	}
	table.Render()
}

func SingleDeployment(cmd *cobra.Command, deployment nvcf.DeploymentResponse) {
	if isJSON(cmd) {
		printJSON(cmd, deployment)
	} else {
		printSingleDeploymentTable(cmd, deployment)
	}
}

func printSingleDeploymentTable(cmd *cobra.Command, deployment nvcf.DeploymentResponse) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{deployment.Deployment.FunctionID, deployment.Deployment.FunctionVersionID, string(deployment.Deployment.FunctionStatus)})
	table.Render()
}

func GPUs(cmd *cobra.Command, clusterGroups []nvcf.ClusterGroupsResponseClusterGroup) {
	if isJSON(cmd) {
		printJSON(cmd, clusterGroups)
	} else {
		printGPUsTable(cmd, clusterGroups)
	}
}

func printGPUsTable(cmd *cobra.Command, clusterGroups []nvcf.ClusterGroupsResponseClusterGroup) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"inst_backend", "inst_gpu_type", "inst_type"})
	table.SetBorder(false)

	for _, clusterGroup := range clusterGroups {
		for _, gpu := range clusterGroup.GPUs {
			for _, instanceType := range gpu.InstanceTypes {
				table.Append([]string{
					clusterGroup.Name, // inst_backend
					gpu.Name,          // inst_gpu_type
					instanceType.Name, // inst_type
				})
			}
		}
	}

	table.Render()
}


func PrintASCIIArt(cmd *cobra.Command) {
	asciiArt := NVIDIA_LOGO_2
	customGreen := color.New(color.FgHiGreen)
	customGreenAsciiArt := customGreen.Sprint(asciiArt)
	fmt.Fprintln(cmd.OutOrStdout(), customGreenAsciiArt)
}

var NVIDIA_LOGO_1 = `
                            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                      @@@@@@@       @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                 @@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
             @@@@@@@@@@@    @@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@@                    
          @@@@@@@@@@     @@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@                    
        @@@@@@@@     @@@@@@@@       @@@@@@@@@    @@@@@@@@@@@@@@@@@@@                    
      @@@@@@@@    @@@@@@@@@@@@@@       @@@@@@@     @@@@@@@@@@@@@@@@@                    
     @@@@@@@    @@@@@@@     @@@@@@      @@@@@@@     @@@@@@@@@@@@@@@@                    
      @@@@@@@   @@@@@@      @@@@@@@   @@@@@@@     @@@@@@@@@@@@@@@@@@                    
       @@@@@@@   @@@@@@     @@@@@@@@@@@@@@@@    @@@@@@@@@@@@@@@@@@@@                    
        @@@@@@@   @@@@@@    @@@@@@@@@@@@@@    @@@@@@@@@@ @@@@@@@@@@@                    
         @@@@@@@   @@@@@@@@ @@@@@@@@@@@@   @@@@@@@@@@       @@@@@@@@                    
          @@@@@@@@   @@@@@@@@@@@@@@@    @@@@@@@@@@@          @@@@@@@                    
            @@@@@@@@    @@@@@      @@@@@@@@@@@@@          @@@@@@@@@@                    
              @@@@@@@@@     @@@@@@@@@@@@@@@@          @@@@@@@@@@@@@@                    
                 @@@@@@@@@  @@@@@@@@@@@           @@@@@@@@@@@@@@@@@@                    
                   @@@@@@@@@@ @@            @@@@@@@@@@@@@@@@@@@@@@@@                    
                       @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                             @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                    
                                                                                        
                                                                                        
                                                                                        
@@@@@@@@@@@@@   @@@@@@     @@@@@@ @@@@@  @@@@@@@@@@@@@    @@@@@@      @@@@@@@@          
@@@@@@@@@@@@@@@ @@@@@@     @@@@@  @@@@@  @@@@@@@@@@@@@@@  @@@@@@      @@@@@@@@          
@@@@@@@@@@@@@@@@ @@@@@@   @@@@@@  @@@@@  @@@@@@@@@@@@@@@  @@@@@@     @@@@@@@@@@         
@@@@@     @@@@@@  @@@@@   @@@@@   @@@@@  @@@@@      @@@@@ @@@@@@    @@@@@ @@@@@@        
@@@@@      @@@@@  @@@@@@ @@@@@    @@@@@  @@@@@      @@@@@ @@@@@@   @@@@@@  @@@@@@       
@@@@@      @@@@@   @@@@@ @@@@@    @@@@@  @@@@@      @@@@@ @@@@@@   @@@@@@@@@@@@@@       
@@@@@      @@@@@   @@@@@@@@@@     @@@@@  @@@@@   @@@@@@@@ @@@@@@  @@@@@@@@@@@@@@@@      
@@@@@      @@@@@    @@@@@@@@@     @@@@@  @@@@@@@@@@@@@@@  @@@@@@ @@@@@@@@@@@@@@@@@@ @@@@
@@@@@      @@@@@     @@@@@@@      @@@@@  @@@@@@@@@@@@@    @@@@@  @@@@@        @@@@@ @@@@
`

var NVIDIA_LOGO_2 = `
                   @@@@@@@@@@@@@@@@@@@@@@@@@@@             
                   @@@@@@@@@@@@@@@@@@@@@@@@@@@             
             @@@@@@@     @@@@@@@@@@@@@@@@@@@@@             
         @@@@@@@@  @@@@@@@   @@@@@@@@@@@@@@@@@             
      @@@@@@    @@@@@@@@@@@@@   @@@@@@@@@@@@@@             
    @@@@@@   @@@@@@@@     @@@@@   @@@@@@@@@@@@             
   @@@@@   @@@@@   @@@@   @@@@@    @@@@@@@@@@@             
    @@@@@@ @@@@@   @@@@@@@@@@@  @@@@@@@@@@@@@@             
     @@@@@@ @@@@@  @@@@@@@@@  @@@@@@@  @@@@@@@             
       @@@@@  @@@@@@@@@@@  @@@@@@@       @@@@@             
         @@@@@   @@@@@@@@@@@@@@       @@@@@@@@             
           @@@@@@@ @@@@@@@        @@@@@@@@@@@@             
              @@@@@@      @@@@@@@@@@@@@@@@@@@@             
                   @@@@@@@@@@@@@@@@@@@@@@@@@@@             
                   @@@@@@@@@@@@@@@@@@@@@@@@@@@             
                                                           
                                                           
@@@@@@@@@ @@@@@   @@@@ @@@ @@@@@@@@@   @@@@    @@@@@       
@@@@@@@@@@ @@@@   @@@@ @@@ @@@@@@@@@@@ @@@@   @@@@@@@      
@@@@   @@@@ @@@@ @@@@  @@@ @@@@    @@@ @@@@  @@@@ @@@@     
@@@@   @@@@ @@@@ @@@   @@@ @@@@    @@@@@@@@  @@@@  @@@@    
@@@@   @@@@  @@@@@@@   @@@ @@@@@@@@@@@ @@@@ @@@@@@@@@@@ @@@
@@@@   @@@@   @@@@@    @@@ @@@@@@@@@@  @@@@@@@@     @@@@@@@
`

var NVCF_LOGO = `
		███╗   ██╗██╗   ██╗ ██████╗███████╗
		████╗  ██║██║   ██║██╔════╝██╔════╝
		██╔██╗ ██║██║   ██║██║     █████╗  
		██║╚██╗██║╚██╗ ██╔╝██║     ██╔══╝  
		██║ ╚████║ ╚████╔╝ ╚██████╗██║     
		╚═╝  ╚═══╝  ╚═══╝   ╚═════╝╚═╝    
`

var NVIDIA_LOGO_3 = `
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                              @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                        @@@@@@@@@@          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                    @@@@@@@@@@@@@@@@@@@@          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                 @@@@@@@@@       @@@@@@@@@@@@        @@@@@@@@@@@@@@@@@@@@@@@@@@@                       
             @@@@@@@@@           @@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@@@                       
           @@@@@@@@@       @@@@@@@       @@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@@                       
         @@@@@@@@      @@@@@@@@@@@           @@@@@@@@       @@@@@@@@@@@@@@@@@@@@                       
      @@@@@@@@@@     @@@@@@@@    @@@@@         @@@@@@@@      @@@@@@@@@@@@@@@@@@@                       
      @@@@@@@@     @@@@@@@@      @@@@@@@       @@@@@@@       @@@@@@@@@@@@@@@@@@@                       
       @@@@@@@@     @@@@@@       @@@@@@@@   @@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@                       
        @@@@@@@@    @@@@@@@      @@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@                       
         @@@@@@@@    @@@@@@@     @@@@@@@@@@@@@@@@      @@@@@@@@@@  @@@@@@@@@@@@@                       
          @@@@@@@@     @@@@@@@   @@@@@@@@@@@@@@     @@@@@@@@@@@       @@@@@@@@@@                       
            @@@@@@@@    @@@@@@@@ @@@@@@@@@@@      @@@@@@@@@@             @@@@@@@                       
             @@@@@@@@      @@@@@@@             @@@@@@@@@@              @@@@@@@@@                       
               @@@@@@@@@      @@@@       @@@@@@@@@@@@@@            @@@@@@@@@@@@@                       
                 @@@@@@@@@@      @@@@@@@@@@@@@@@@@@            @@@@@@@@@@@@@@@@@                       
                    @@@@@@@@@@@  @@@@@@@@@@@@@             @@@@@@@@@@@@@@@@@@@@@                       
                       @@@@@@@@@@@                   @@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                           @@@@@@@           @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       
                                                                                                       
                                                                                                       
                                                                                                       
                                                                                                       
@@@@@@@@@@@@       @@@@@@        @@@@@  @@@@@@   @@@@@@@@@           @@@@@         @@@@@@@             
@@@@@@@@@@@@@@@@   @@@@@@       @@@@@@  @@@@@@  @@@@@@@@@@@@@@@@    @@@@@@        @@@@@@@@@            
@@@@@@@@@@@@@@@@@@  @@@@@@     @@@@@@@  @@@@@@  @@@@@@@@@@@@@@@@@   @@@@@@       @@@@@@@@@@@           
@@@@@@     @@@@@@@  @@@@@@@    @@@@@@   @@@@@@  @@@@@@     @@@@@@@  @@@@@@       @@@@@@@@@@@           
@@@@@@       @@@@@   @@@@@@   @@@@@@@   @@@@@@  @@@@@@       @@@@@  @@@@@@      @@@@@@ @@@@@@          
@@@@@@       @@@@@@  @@@@@@@  @@@@@@    @@@@@@  @@@@@@       @@@@@@ @@@@@@     @@@@@@   @@@@@@         
@@@@@@       @@@@@@   @@@@@@ @@@@@@     @@@@@@  @@@@@@       @@@@@  @@@@@@    @@@@@@     @@@@@@        
@@@@@@       @@@@@@    @@@@@@@@@@@@     @@@@@@  @@@@@@      @@@@@@  @@@@@@    @@@@@@@@@@@@@@@@@        
@@@@@@       @@@@@@    @@@@@@@@@@@      @@@@@@  @@@@@@@@@@@@@@@@@@  @@@@@@   @@@@@@@@@@@@@@@@@@@   @@@@
@@@@@@       @@@@@@     @@@@@@@@@@      @@@@@@  @@@@@@@@@@@@@@@@@   @@@@@@  @@@@@@        @@@@@@@ @@@@@
@@@@@@       @@@@@@     @@@@@@@@@       @@@@@@  @@@@@@@@@@@@@@      @@@@@@  @@@@@          @@@@@@ @@@@@
`

type TypeOptions struct {
	Speed       time.Duration
	Skippable   bool
	Writer      io.Writer
	StopChannel chan struct{}
}

var defaultOptions = TypeOptions{
	Speed:       27 * time.Millisecond,
	Skippable:   true,
	Writer:      os.Stdout,
	StopChannel: nil,
}

func Type(s string, opts ...TypeOptions) {
	options := mergeOptions(defaultOptions, opts...)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		typeText(s, options)
	}()

	if options.Skippable {
		go func() {
			fmt.Scanln() // Wait for Enter key
			if options.StopChannel != nil {
				close(options.StopChannel)
			}
		}()
	}

	wg.Wait()
}

func mergeOptions(defaultOpts TypeOptions, opts ...TypeOptions) TypeOptions {
	if len(opts) == 0 {
		return defaultOpts
	}
	userOpts := opts[0]
	if userOpts.Speed == 0 {
		userOpts.Speed = defaultOpts.Speed
	}
	if userOpts.Writer == nil {
		userOpts.Writer = defaultOpts.Writer
	}
	if userOpts.StopChannel == nil {
		userOpts.StopChannel = defaultOpts.StopChannel
	}
	return userOpts
}

func typeText(s string, options TypeOptions) {
	for _, char := range s {
		select {
		case <-options.StopChannel:
			fmt.Fprint(options.Writer, s[len(s)-len(string(char)):])
			return
		default:
			fmt.Fprintf(options.Writer, "%c", char)
			time.Sleep(options.Speed)
		}
	}
}

// TypeWithColor types text with the specified color
func TypeWithColor(s string, c *color.Color, opts ...TypeOptions) {
	options := mergeOptions(defaultOptions, opts...)
	coloredString := c.SprintFunc()(s)
	Type(coloredString, options)
}

// Update existing output functions to use Type
func Success(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		TypeWithColor(message+"\n", color.New(color.FgGreen))
	}
}

func Info(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		TypeWithColor(message+"\n", color.New(color.FgBlue))
	}
}


// Example usage:
func ExampleTyping() {
	// Basic usage
	Type("Hello, World! Messages can go here.\n")

	// Custom options
	Type("This types at default speed and can't be skipped\n", TypeOptions{
		Skippable: false,
	})

	// With color
	TypeWithColor("This is a green message\n", color.New(color.FgGreen))

	// With stop channel
	stopChan := make(chan struct{})
	Type("Press Enter to skip this message...\n", TypeOptions{
		StopChannel: stopChan,
	})
}