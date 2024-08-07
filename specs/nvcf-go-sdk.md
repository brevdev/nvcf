Here is the API documentation for the provided Go codebase:

<api_documentation>
  <package name="nvcf">
    <structs>
      <struct>
        <name>Client</name>
        <fields>
          <field>Options []option.RequestOption</field>
          <field>FunctionManagement *FunctionManagementService</field>
          <field>FunctionDeployment *FunctionDeploymentService</field>
          <field>FunctionInvocation *FunctionInvocationService</field>
          <field>EnvelopeFunctionInvocation *EnvelopeFunctionInvocationService</field>
          <field>Functions *FunctionService</field>
          <field>Authorizations *AuthorizationService</field>
          <field>Assets *AssetService</field>
          <field>AssetManagement *AssetManagementService</field>
          <field>AuthorizedAccounts *AuthorizedAccountService</field>
          <field>QueueDetails *QueueDetailService</field>
          <field>Exec *ExecService</field>
          <field>ClusterGroups *ClusterGroupService</field>
          <field>Clients *ClientService</field>
        </fields>
        <methods>
          <method>
            <name>NewClient</name>
            <signature>func NewClient(opts ...option.RequestOption) *Client</signature>
            <description>NewClient generates a new client with the default options and any additional options provided</description>
          </method>
          <method>
            <name>Execute</name>
            <signature>func (r *Client) Execute(ctx context.Context, method string, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Execute makes a request with the given context, method, URL, request params, response, and request options</description>
          </method>
          <method>
            <name>Get</name>
            <signature>func (r *Client) Get(ctx context.Context, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Get makes a GET request with the given URL, params, and optionally deserializes to a response</description>
          </method>
          <method>
            <name>Post</name>
            <signature>func (r *Client) Post(ctx context.Context, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Post makes a POST request with the given URL, params, and optionally deserializes to a response</description>
          </method>
          <method>
            <name>Put</name>
            <signature>func (r *Client) Put(ctx context.Context, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Put makes a PUT request with the given URL, params, and optionally deserializes to a response</description>
          </method>
          <method>
            <name>Patch</name>
            <signature>func (r *Client) Patch(ctx context.Context, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Patch makes a PATCH request with the given URL, params, and optionally deserializes to a response</description>
          </method>
          <method>
            <name>Delete</name>
            <signature>func (r *Client) Delete(ctx context.Context, path string, params interface{}, res interface{}, opts ...option.RequestOption) error</signature>
            <description>Delete makes a DELETE request with the given URL, params, and optionally deserializes to a response</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>AssetResponse</name>
        <fields>
          <field>Asset AssetResponseAsset</field>
          <field>JSON assetResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *AssetResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>AssetResponseAsset</name>
        <fields>
          <field>AssetID string</field>
          <field>ContentType string</field>
          <field>Description string</field>
          <field>JSON assetResponseAssetJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *AssetResponseAsset) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListAssetsResponse</name>
        <fields>
          <field>Assets []ListAssetsResponseAsset</field>
          <field>JSON listAssetsResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListAssetsResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListAssetsResponseAsset</name>
        <fields>
          <field>AssetID string</field>
          <field>ContentType string</field>
          <field>Description string</field>
          <field>JSON listAssetsResponseAssetJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListAssetsResponseAsset) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateAssetResponse</name>
        <fields>
          <field>AssetID string</field>
          <field>ContentType string</field>
          <field>Description string</field>
          <field>UploadURL string</field>
          <field>JSON createAssetResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateAssetResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ClusterGroupsResponse</name>
        <fields>
          <field>ClusterGroups []ClusterGroupsResponseClusterGroup</field>
          <field>JSON clusterGroupsResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ClusterGroupsResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ClusterGroupsResponseClusterGroup</name>
        <fields>
          <field>ID string</field>
          <field>AuthorizedNcaIDs []string</field>
          <field>Clusters []ClusterGroupsResponseClusterGroupsCluster</field>
          <field>GPUs []ClusterGroupsResponseClusterGroupsGPU</field>
          <field>Name string</field>
          <field>NcaID string</field>
          <field>JSON clusterGroupsResponseClusterGroupJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ClusterGroupsResponseClusterGroup) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ClusterGroupsResponseClusterGroupsCluster</name>
        <fields>
          <field>ID string</field>
          <field>K8sVersion string</field>
          <field>Name string</field>
          <field>JSON clusterGroupsResponseClusterGroupsClusterJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ClusterGroupsResponseClusterGroupsCluster) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ClusterGroupsResponseClusterGroupsGPU</name>
        <fields>
          <field>InstanceTypes []ClusterGroupsResponseClusterGroupsGPUsInstanceType</field>
          <field>Name string</field>
          <field>JSON clusterGroupsResponseClusterGroupsGPUJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ClusterGroupsResponseClusterGroupsGPU) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ClusterGroupsResponseClusterGroupsGPUsInstanceType</name>
        <fields>
          <field>Default bool</field>
          <field>Description string</field>
          <field>Name string</field>
          <field>JSON clusterGroupsResponseClusterGroupsGPUsInstanceTypeJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ClusterGroupsResponseClusterGroupsGPUsInstanceType) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponse</name>
        <fields>
          <field>Function CreateFunctionResponseFunction</field>
          <field>JSON createFunctionResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunction</name>
        <fields>
          <field>ID string</field>
          <field>CreatedAt time.Time</field>
          <field>FunctionType CreateFunctionResponseFunctionFunctionType</field>
          <field>HealthUri string</field>
          <field>Name string</field>
          <field>NcaID string</field>
          <field>Status CreateFunctionResponseFunctionStatus</field>
          <field>VersionID string</field>
          <field>ActiveInstances []CreateFunctionResponseFunctionActiveInstance</field>
          <field>APIBodyFormat CreateFunctionResponseFunctionAPIBodyFormat</field>
          <field>ContainerArgs string</field>
          <field>ContainerEnvironment []CreateFunctionResponseFunctionContainerEnvironment</field>
          <field>ContainerImage string</field>
          <field>Description string</field>
          <field>Health CreateFunctionResponseFunctionHealth</field>
          <field>HelmChart string</field>
          <field>HelmChartServiceName string</field>
          <field>InferencePort int64</field>
          <field>InferenceURL string</field>
          <field>Models []CreateFunctionResponseFunctionModel</field>
          <field>OwnedByDifferentAccount bool</field>
          <field>Resources []CreateFunctionResponseFunctionResource</field>
          <field>Tags []string</field>
          <field>JSON createFunctionResponseFunctionJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunction) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunctionActiveInstance</name>
        <fields>
          <field>Backend string</field>
          <field>FunctionID string</field>
          <field>FunctionVersionID string</field>
          <field>GPU string</field>
          <field>InstanceCreatedAt time.Time</field>
          <field>InstanceID string</field>
          <field>InstanceStatus CreateFunctionResponseFunctionActiveInstancesInstanceStatus</field>
          <field>InstanceType string</field>
          <field>InstanceUpdatedAt time.Time</field>
          <field>Location string</field>
          <field>NcaID string</field>
          <field>SisRequestID string</field>
          <field>JSON createFunctionResponseFunctionActiveInstanceJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunctionActiveInstance) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunctionContainerEnvironment</name>
        <fields>
          <field>Key string</field>
          <field>Value string</field>
          <field>JSON createFunctionResponseFunctionContainerEnvironmentJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunctionContainerEnvironment) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunctionHealth</name>
        <fields>
          <field>ExpectedStatusCode int64</field>
          <field>Port int64</field>
          <field>Protocol CreateFunctionResponseFunctionHealthProtocol</field>
          <field>Timeout string</field>
          <field>Uri string</field>
          <field>JSON createFunctionResponseFunctionHealthJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunctionHealth) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunctionModel</name>
        <fields>
          <field>Name string</field>
          <field>Uri string</field>
          <field>Version string</field>
          <field>JSON createFunctionResponseFunctionModelJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunctionModel) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>CreateFunctionResponseFunctionResource</name>
        <fields>
          <field>Name string</field>
          <field>Uri string</field>
          <field>Version string</field>
          <field>JSON createFunctionResponseFunctionResourceJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *CreateFunctionResponseFunctionResource) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponse</name>
        <fields>
          <field>Functions []ListFunctionsResponseFunction</field>
          <field>JSON listFunctionsResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunction</name>
        <fields>
          <field>ID string</field>
          <field>CreatedAt time.Time</field>
          <field>FunctionType ListFunctionsResponseFunctionsFunctionType</field>
          <field>HealthUri string</field>
          <field>Name string</field>
          <field>NcaID string</field>
          <field>Status ListFunctionsResponseFunctionsStatus</field>
          <field>VersionID string</field>
          <field>ActiveInstances []ListFunctionsResponseFunctionsActiveInstance</field>
          <field>APIBodyFormat ListFunctionsResponseFunctionsAPIBodyFormat</field>
          <field>ContainerArgs string</field>
          <field>ContainerEnvironment []ListFunctionsResponseFunctionsContainerEnvironment</field>
          <field>ContainerImage string</field>
          <field>Description string</field>
          <field>Health ListFunctionsResponseFunctionsHealth</field>
          <field>HelmChart string</field>
          <field>HelmChartServiceName string</field>
          <field>InferencePort int64</field>
          <field>InferenceURL string</field>
          <field>Models []ListFunctionsResponseFunctionsModel</field>
          <field>OwnedByDifferentAccount bool</field>
          <field>Resources []ListFunctionsResponseFunctionsResource</field>
          <field>Tags []string</field>
          <field>JSON listFunctionsResponseFunctionJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunction) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunctionsActiveInstance</name>
        <fields>
          <field>Backend string</field>
          <field>FunctionID string</field>
          <field>FunctionVersionID string</field>
          <field>GPU string</field>
          <field>InstanceCreatedAt time.Time</field>
          <field>InstanceID string</field>
          <field>InstanceStatus ListFunctionsResponseFunctionsActiveInstancesInstanceStatus</field>
          <field>InstanceType string</field>
          <field>InstanceUpdatedAt time.Time</field>
          <field>Location string</field>
          <field>NcaID string</field>
          <field>SisRequestID string</field>
          <field>JSON listFunctionsResponseFunctionsActiveInstanceJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunctionsActiveInstance) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunctionsContainerEnvironment</name>
        <fields>
          <field>Key string</field>
          <field>Value string</field>
          <field>JSON listFunctionsResponseFunctionsContainerEnvironmentJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunctionsContainerEnvironment) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunctionsHealth</name>
        <fields>
          <field>ExpectedStatusCode int64</field>
          <field>Port int64</field>
          <field>Protocol ListFunctionsResponseFunctionsHealthProtocol</field>
          <field>Timeout string</field>
          <field>Uri string</field>
          <field>JSON listFunctionsResponseFunctionsHealthJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunctionsHealth) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunctionsModel</name>
        <fields>
          <field>Name string</field>
          <field>Uri string</field>
          <field>Version string</field>
          <field>JSON listFunctionsResponseFunctionsModelJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunctionsModel) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionsResponseFunctionsResource</name>
        <fields>
          <field>Name string</field>
          <field>Uri string</field>
          <field>Version string</field>
          <field>JSON listFunctionsResponseFunctionsResourceJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionsResponseFunctionsResource) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListFunctionIDsResponse</name>
        <fields>
          <field>FunctionIDs []string</field>
          <field>JSON listFunctionIDsResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListFunctionIDsResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListAuthorizedPartiesResponse</name>
        <fields>
          <field>Functions []ListAuthorizedPartiesResponseFunction</field>
          <field>JSON listAuthorizedPartiesResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListAuthorizedPartiesResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListAuthorizedPartiesResponseFunction</name>
        <fields>
          <field>ID string</field>
          <field>NcaID string</field>
          <field>AuthorizedParties []ListAuthorizedPartiesResponseFunctionsAuthorizedParty</field>
          <field>VersionID string</field>
          <field>JSON listAuthorizedPartiesResponseFunctionJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListAuthorizedPartiesResponseFunction) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>ListAuthorizedPartiesResponseFunctionsAuthorizedParty</name>
        <fields>
          <field>NcaID string</field>
          <field>ClientID string</field>
          <field>JSON listAuthorizedPartiesResponseFunctionsAuthorizedPartyJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *ListAuthorizedPartiesResponseFunctionsAuthorizedParty) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>DeploymentResponse</name>
        <fields>
          <field>Deployment DeploymentResponseDeployment</field>
          <field>JSON deploymentResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *DeploymentResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>DeploymentResponseDeployment</name>
        <fields>
          <field>DeploymentSpecifications []DeploymentResponseDeploymentDeploymentSpecification</field>
          <field>FunctionID string</field>
          <field>FunctionStatus DeploymentResponseDeploymentFunctionStatus</field>
          <field>FunctionVersionID string</field>
          <field>NcaID string</field>
          <field>HealthInfo []DeploymentResponseDeploymentHealthInfo</field>
          <field>RequestQueueURL string</field>
          <field>JSON deploymentResponseDeploymentJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *DeploymentResponseDeployment) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>DeploymentResponseDeploymentDeploymentSpecification</name>
        <fields>
          <field>GPU string</field>
          <field>InstanceType string</field>
          <field>MaxInstances int64</field>
          <field>MinInstances int64</field>
          <field>AvailabilityZones []string</field>
          <field>Backend string</field>
          <field>Configuration interface{}</field>
          <field>MaxRequestConcurrency int64</field>
          <field>JSON deploymentResponseDeploymentDeploymentSpecificationJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *DeploymentResponseDeploymentDeploymentSpecification) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>DeploymentResponseDeploymentHealthInfo</name>
        <fields>
          <field>Backend string</field>
          <field>Error string</field>
          <field>GPU string</field>
          <field>InstanceType string</field>
          <field>SisRequestID string</field>
          <field>JSON deploymentResponseDeploymentHealthInfoJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *DeploymentResponseDeploymentHealthInfo) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>FunctionInvocationFunctionInvokeResponse</name>
        <fields>
          <field>Char string</field>
          <field>Direct bool</field>
          <field>Double float64</field>
          <field>Float float64</field>
          <field>Int int64</field>
          <field>Long int64</field>
          <field>ReadOnly bool</field>
          <field>Short int64</field>
          <field>JSON functionInvocationFunctionInvokeResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *FunctionInvocationFunctionInvokeResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>FunctionInvocationFunctionVersionInvokeResponse</name>
        <fields>
          <field>Char string</field>
          <field>Direct bool</field>
          <field>Double float64</field>
          <field>Float float64</field>
          <field>Int int64</field>
          <field>Long int64</field>
          <field>ReadOnly bool</field>
          <field>Short int64</field>
          <field>JSON functionInvocationFunctionVersionInvokeResponseJSON</field>
        </fields>
        <methods>
          <method>
            <name>UnmarshalJSON</name>
            <signature>func (r *FunctionInvocationFunctionVersionInvokeResponse) UnmarshalJSON(data []byte) (err error)</signature>
            <description>UnmarshalJSON implements the json.Unmarshaler interface</description>
          </method>
        </methods>
      </struct>
      
      <struct>
        <name>FunctionInvocationStatusGetResponse</name>
        <fields>
          <field>Char string</field>
          <field>Direct bool</field>
          <field>Double float64</field>
          <field>Float float64</field>
          <field>Int int64</field>