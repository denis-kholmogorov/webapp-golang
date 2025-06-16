package domain

type Settings struct {
	Id                     string   `json:"id,omitempty"`
	Uid                    string   `json:"uid,omitempty"`
	EnablePost             bool     `json:"enablePost"`
	EnablePostComment      bool     `json:"enablePostComment"`
	EnableCommentComment   bool     `json:"enableCommentComment"`
	EnableMessage          bool     `json:"enableMessage"`
	EnableFriendRequest    bool     `json:"enableFriendRequest"`
	EnableFriendBirthday   bool     `json:"enableFriendBirthday"`
	EnableSendEmailMessage bool     `json:"enableSendEmailMessage"`
	IsDeleted              bool     `json:"isDeleted"`
	DType                  []string `json:"dgraph.type,omitempty"`
}

type SettingsList struct {
	List []Settings `json:"settingsList,omitempty"`
}

func NewSettings() Settings {
	return Settings{
		Uid:                    "_:settings",
		EnablePost:             true,
		EnablePostComment:      true,
		EnableCommentComment:   true,
		EnableMessage:          true,
		EnableFriendRequest:    true,
		EnableFriendBirthday:   true,
		EnableSendEmailMessage: true,
		IsDeleted:              true,
		DType:                  []string{"Settings"},
	}
}
