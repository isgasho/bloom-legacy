use serde::{Serialize, Deserialize};
use uuid::Uuid;


#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct StartUploadSessionBody {
    pub file_name: String,
    pub parent_id: Option<Uuid>,
}


#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct StartUploadSessionResponse {
    pub id: Uuid,
    pub presigned_url: String,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct CompleteUploadSessionBody {
    pub upload_session_id: Uuid,
}
