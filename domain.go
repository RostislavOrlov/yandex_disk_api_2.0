package main

import "time"

type User struct {
	Id     int
	Name   string
	Token  string
	ChatId int
	Files  []File
}

type UserToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type Item struct {
	Name       string                 `json:"Name"`
	Exif       map[string]interface{} `json:"exif"`
	Created    string                 `json:"created"`
	ResourceId string                 `json:"resource_id"`
	Modified   string                 `json:"modified"`
	Path       string                 `json:"path"`
	CommentIds struct {
		PrivateResource string `json:"private_resource"`
		PublicResource  string `json:"public_resource"`
	} `json:"comment_ids"`
	Type     string `json:"type"`
	Revision int    `json:"revision"`
}

type FetchFoldersAndViewsResponse struct {
	Embedded struct {
		Sort   string `json:"sort"`
		Items  []Item `json:"items"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Path   string `json:"path"`
		Total  int    `json:"total"`
	} `json:"_embedded"`
	NameDisk       string                 `json:"Name"`
	ExifDisk       map[string]interface{} `json:"exif"`
	ResourceIDDisk string                 `json:"resource_id"`
	CreatedDisk    string                 `json:"created"`
	ModifiedDisk   string                 `json:"modified"`
	PathDisk       string                 `json:"path"`
	CommentIDsDisk struct {
		//PrivateResourceDisk string `json:"private_resource"`
		//PublicResourceDisk  string `json:"public_resource"`
	} `json:"comment_ids"`
	TypeDisk     string  `json:"type"`
	RevisionDisk float64 `json:"revision"`
}

type DownloadFileResponse struct {
	HRef      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

type ShowInfoResponse struct {
	MaxFileSize     int    `json:"max_file_size"`
	PaidMaxFileSize int    `json:"paid_max_file_size"`
	TotalSpace      int    `json:"total_space"`
	RegTime         string `json:"reg_time"`
	TrashSize       int    `json:"trash_size"`
	IsPaid          bool   `json:"is_paid"`
	UsedSpace       int    `json:"used_space"`
}

type ShowInfo struct {
	TotalSpace string `json:"total_space"`
	UsedSpace  string `json:"used_space"`
	TrashSize  string `json:"trash_size"`
}

type File struct {
	Name string
	Path string
}

type FilesResponse struct {
	Items []struct {
		Name       string    `json:"name"`
		Created    time.Time `json:"created"`
		Size       int       `json:"size"`
		CommentIds struct {
			PrivateResource string `json:"private_resource"`
			PublicResource  string `json:"public_resource"`
		} `json:"comment_ids"`
		Sizes []struct {
			Url  string `json:"url"`
			Name string `json:"name"`
		} `json:"sizes,omitempty"`
		File      string `json:"file"`
		MediaType string `json:"media_type"`
		Preview   string `json:"preview,omitempty"`
		Path      string `json:"path"`
		Sha256    string `json:"sha256"`
		Type      string `json:"type"`
		Md5       string `json:"md5"`
		Revision  int64  `json:"revision"`
	} `json:"items"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TelegramMessageResponse struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		MessageId int `json:"message_id"`
		From      struct {
			Id           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type ListKMSkeysResponse struct {
	Keys []struct {
		PrimaryVersion struct {
			Primary     bool      `json:"primary"`
			HostedByHsm bool      `json:"hostedByHsm"`
			Id          string    `json:"id"`
			KeyId       string    `json:"keyId"`
			Status      string    `json:"status"`
			Algorithm   string    `json:"algorithm"`
			CreatedAt   time.Time `json:"createdAt"`
		} `json:"primaryVersion"`
		DeletionProtection bool      `json:"deletionProtection"`
		Id                 string    `json:"id"`
		FolderId           string    `json:"folderId"`
		CreatedAt          time.Time `json:"createdAt"`
		Status             string    `json:"status"`
		DefaultAlgorithm   string    `json:"defaultAlgorithm"`
		Name               string    `json:"name,omitempty"`
	} `json:"keys"`
}

type CreateKMSkeyResponse struct {
	Done     bool `json:"done"`
	Metadata struct {
		Type             string `json:"@type"`
		KeyId            string `json:"keyId"`
		PrimaryVersionId string `json:"primaryVersionId"`
	} `json:"metadata"`
	Response struct {
		Type           string `json:"@type"`
		PrimaryVersion struct {
			Primary     bool      `json:"primary"`
			HostedByHsm bool      `json:"hostedByHsm"`
			Id          string    `json:"id"`
			KeyId       string    `json:"keyId"`
			Status      string    `json:"status"`
			Algorithm   string    `json:"algorithm"`
			CreatedAt   time.Time `json:"createdAt"`
		} `json:"primaryVersion"`
		DeletionProtection bool      `json:"deletionProtection"`
		Id                 string    `json:"id"`
		FolderId           string    `json:"folderId"`
		CreatedAt          time.Time `json:"createdAt"`
		Name               string    `json:"name"`
		Status             string    `json:"status"`
		DefaultAlgorithm   string    `json:"defaultAlgorithm"`
	} `json:"response"`
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

type CreateLockboxSecretResponse struct {
	Done     bool `json:"done"`
	Metadata struct {
		Type     string `json:"@type"`
		SecretId string `json:"secretId"`
	} `json:"metadata"`
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

type GetLockboxPayloadResponse struct {
	Entries []struct {
		Key       string `json:"key"`
		TextValue string `json:"textValue"`
	} `json:"entries"`
	VersionId string `json:"versionId"`
}

type ListLockboxSecretsResponse struct {
	Secrets []struct {
		CurrentVersion struct {
			PayloadEntryKeys []string  `json:"payloadEntryKeys"`
			Id               string    `json:"id"`
			SecretId         string    `json:"secretId"`
			CreatedAt        time.Time `json:"createdAt"`
			Status           string    `json:"status"`
		} `json:"currentVersion"`
		DeletionProtection bool      `json:"deletionProtection"`
		Id                 string    `json:"id"`
		FolderId           string    `json:"folderId"`
		CreatedAt          time.Time `json:"createdAt"`
		Name               string    `json:"name"`
		KmsKeyId           string    `json:"kmsKeyId"`
		Status             string    `json:"status"`
	} `json:"secrets"`
}

type GetIAMtokenResponse struct {
	IamToken  string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}
