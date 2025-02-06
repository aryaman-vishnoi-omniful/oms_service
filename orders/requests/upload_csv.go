package requests

type CSVUploadRequest struct {
	FilePath string `json:"file_path"`
	TenantID   uint64 `json:"tenant_id"` 
	HubID	uint64 `json:"hub_id"`
	

}