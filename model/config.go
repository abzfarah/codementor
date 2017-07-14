// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.


package model

import (
	"encoding/json"
	"io"
	"net/url"
)

const (
	CONN_SECURITY_NONE     = ""
	CONN_SECURITY_PLAIN    = "PLAIN"
	CONN_SECURITY_TLS      = "TLS"
	CONN_SECURITY_STARTTLS = "STARTTLS"

	IMAGE_DRIVER_LOCAL = "local"
	IMAGE_DRIVER_S3    = "amazons3"

	DATABASE_DRIVER_MYSQL    = "mysql"
	DATABASE_DRIVER_POSTGRES = "postgres"

	PASSWORD_MAXIMUM_LENGTH = 64
	PASSWORD_MINIMUM_LENGTH = 5

	SERVICE_GITLAB    = "gitlab"
	SERVICE_GOOGLE    = "google"
	SERVICE_OFFICE365 = "office365"

	WEBSERVER_MODE_REGULAR  = "regular"
	WEBSERVER_MODE_GZIP     = "gzip"
	WEBSERVER_MODE_DISABLED = "disabled"

	GENERIC_NOTIFICATION = "generic"
	FULL_NOTIFICATION    = "full"

	DIRECT_MESSAGE_ANY  = "any"
	DIRECT_MESSAGE_TEAM = "team"

	PERMISSIONS_ALL           = "all"
	PERMISSIONS_CHANNEL_ADMIN = "channel_admin"
	PERMISSIONS_TEAM_ADMIN    = "team_admin"
	PERMISSIONS_SYSTEM_ADMIN  = "system_admin"

	FAKE_SETTING = "********************************"

	RESTRICT_EMOJI_CREATION_ALL          = "all"
	RESTRICT_EMOJI_CREATION_ADMIN        = "admin"
	RESTRICT_EMOJI_CREATION_SYSTEM_ADMIN = "system_admin"

	PERMISSIONS_DELETE_POST_ALL          = "all"
	PERMISSIONS_DELETE_POST_TEAM_ADMIN   = "team_admin"
	PERMISSIONS_DELETE_POST_SYSTEM_ADMIN = "system_admin"

	ALLOW_EDIT_POST_ALWAYS     = "always"
	ALLOW_EDIT_POST_NEVER      = "never"
	ALLOW_EDIT_POST_TIME_LIMIT = "time_limit"

	EMAIL_BATCHING_BUFFER_SIZE = 256
	EMAIL_BATCHING_INTERVAL    = 30

	SITENAME_MAX_LENGTH = 30

	SERVICE_SETTINGS_DEFAULT_SITE_URL        = ""
	SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE   = ""
	SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE    = ""
	SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT    = 300
	SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT   = 300
	SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM = ""

	TEAM_SETTINGS_DEFAULT_CUSTOM_BRAND_TEXT        = ""
	TEAM_SETTINGS_DEFAULT_CUSTOM_DESCRIPTION_TEXT  = ""
	TEAM_SETTINGS_DEFAULT_USER_STATUS_AWAY_TIMEOUT = 300

	EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION = ""

	SUPPORT_SETTINGS_DEFAULT_TERMS_OF_SERVICE_LINK = "https://about.mattermost.com/default-terms/"
	SUPPORT_SETTINGS_DEFAULT_PRIVACY_POLICY_LINK   = "https://about.mattermost.com/default-privacy-policy/"
	SUPPORT_SETTINGS_DEFAULT_ABOUT_LINK            = "https://about.mattermost.com/default-about/"
	SUPPORT_SETTINGS_DEFAULT_HELP_LINK             = "https://about.mattermost.com/default-help/"
	SUPPORT_SETTINGS_DEFAULT_REPORT_A_PROBLEM_LINK = "https://about.mattermost.com/default-report-a-problem/"
	SUPPORT_SETTINGS_DEFAULT_SUPPORT_EMAIL         = "feedback@mattermost.com"

	LDAP_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE = ""
	LDAP_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE  = ""
	LDAP_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE      = ""
	LDAP_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE   = ""
	LDAP_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE   = ""
	LDAP_SETTINGS_DEFAULT_ID_ATTRIBUTE         = ""
	LDAP_SETTINGS_DEFAULT_POSITION_ATTRIBUTE   = ""
	LDAP_SETTINGS_DEFAULT_LOGIN_FIELD_NAME     = ""

	SAML_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE = ""
	SAML_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE  = ""
	SAML_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE      = ""
	SAML_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE   = ""
	SAML_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE   = ""
	SAML_SETTINGS_DEFAULT_LOCALE_ATTRIBUTE     = ""
	SAML_SETTINGS_DEFAULT_POSITION_ATTRIBUTE   = ""

	NATIVEAPP_SETTINGS_DEFAULT_APP_DOWNLOAD_LINK         = "https://about.mattermost.com/downloads/"
	NATIVEAPP_SETTINGS_DEFAULT_ANDROID_APP_DOWNLOAD_LINK = "https://about.mattermost.com/mattermost-android-app/"
	NATIVEAPP_SETTINGS_DEFAULT_IOS_APP_DOWNLOAD_LINK     = "https://about.mattermost.com/mattermost-ios-app/"

	WEBRTC_SETTINGS_DEFAULT_STUN_URI = ""
	WEBRTC_SETTINGS_DEFAULT_TURN_URI = ""

	ANALYTICS_SETTINGS_DEFAULT_MAX_USERS_FOR_STATISTICS = 2500
)

type ServiceSettings struct {
	SiteURL                                  *string
	LicenseFileLocation                      *string
	ListenAddress                            string
	ConnectionSecurity                       *string
	TLSCertFile                              *string
	TLSKeyFile                               *string
	UseLetsEncrypt                           *bool
	LetsEncryptCertificateCacheFile          *string
	Forward80To443                           *bool
	ReadTimeout                              *int
	WriteTimeout                             *int
	MaximumLoginAttempts                     int
	GoogleDeveloperKey                       string
	EnableOAuthServiceProvider               bool
	EnableIncomingWebhooks                   bool
	EnableOutgoingWebhooks                   bool
	EnableCommands                           *bool
	EnableOnlyAdminIntegrations              *bool
	EnablePostUsernameOverride               bool
	EnablePostIconOverride                   bool
	EnableLinkPreviews                       *bool
	EnableTesting                            bool
	EnableDeveloper                          *bool
	EnableSecurityFixAlert                   *bool
	EnableInsecureOutgoingConnections        *bool
	EnableMultifactorAuthentication          *bool
	EnforceMultifactorAuthentication         *bool
	AllowCorsFrom                            *string
	SessionLengthWebInDays                   *int
	SessionLengthMobileInDays                *int
	SessionLengthSSOInDays                   *int
	SessionCacheInMinutes                    *int
	WebsocketSecurePort                      *int
	WebsocketPort                            *int
	WebserverMode                            *string
	EnableCustomEmoji                        *bool
	RestrictCustomEmojiCreation              *string
	RestrictPostDelete                       *string
	AllowEditPost                            *string
	PostEditTimeLimit                        *int
	TimeBetweenUserTypingUpdatesMilliseconds *int64
	EnablePostSearch                         *bool
	EnableUserTypingMessages                 *bool
	EnableUserStatuses                       *bool
	ClusterLogTimeoutMilliseconds            *int
}

type ClusterSettings struct {
	Enable                 *bool
	InterNodeListenAddress *string
	InterNodeUrls          []string
}

type MetricsSettings struct {
	Enable           *bool
	BlockProfileRate *int
	ListenAddress    *string
}

type AnalyticsSettings struct {
	MaxUsersForStatistics *int
}

type SSOSettings struct {
	Enable          bool
	Secret          string
	Id              string
	Scope           string
	AuthEndpoint    string
	TokenEndpoint   string
	UserApiEndpoint string
}

type SqlSettings struct {
	DriverName               string
	DataSource               string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	MaxIdleConns             int
	MaxOpenConns             int
	Trace                    bool
	AtRestEncryptKey         string
}

type LogSettings struct {
	EnableConsole          bool
	ConsoleLevel           string
	EnableFile             bool
	FileLevel              string
	FileFormat             string
	FileLocation           string
	EnableWebhookDebugging bool
	EnableDiagnostics      *bool
}

type PasswordSettings struct {
	MinimumLength *int
	Lowercase     *bool
	Number        *bool
	Uppercase     *bool
	Symbol        *bool
}

type FileSettings struct {
	MaxFileSize             *int64
	DriverName              string
	Directory               string
	EnablePublicLink        bool
	PublicLinkSalt          *string
	ThumbnailWidth          int
	ThumbnailHeight         int
	PreviewWidth            int
	PreviewHeight           int
	ProfileWidth            int
	ProfileHeight           int
	InitialFont             string
	AmazonS3AccessKeyId     string
	AmazonS3SecretAccessKey string
	AmazonS3Bucket          string
	AmazonS3Region          string
	AmazonS3Endpoint        string
	AmazonS3SSL             *bool
}

type EmailSettings struct {
	EnableSignUpWithEmail             bool
	EnableSignInWithEmail             *bool
	EnableSignInWithUsername          *bool
	SendEmailNotifications            bool
	RequireEmailVerification          bool
	FeedbackName                      string
	FeedbackEmail                     string
	FeedbackOrganization              *string
	SMTPUsername                      string
	SMTPPassword                      string
	SMTPServer                        string
	SMTPPort                          string
	ConnectionSecurity                string
	InviteSalt                        string
	PasswordResetSalt                 string
	SendPushNotifications             *bool
	PushNotificationServer            *string
	PushNotificationContents          *string
	EnableEmailBatching               *bool
	EmailBatchingBufferSize           *int
	EmailBatchingInterval             *int
	SkipServerCertificateVerification *bool
}

type RateLimitSettings struct {
	Enable           *bool
	PerSec           int
	MaxBurst         *int
	MemoryStoreSize  int
	VaryByRemoteAddr bool
	VaryByHeader     string
}

type PrivacySettings struct {
	ShowEmailAddress bool
	ShowFullName     bool
}

type SupportSettings struct {
	TermsOfServiceLink *string
	PrivacyPolicyLink  *string
	AboutLink          *string
	HelpLink           *string
	ReportAProblemLink *string
	SupportEmail       *string
}

type TeamSettings struct {
	SiteName                            string
	MaxUsersPerTeam                     int
	EnableTeamCreation                  bool
	EnableUserCreation                  bool
	EnableOpenServer                    *bool
	RestrictCreationToDomains           string
	EnableCustomBrand                   *bool
	CustomBrandText                     *string
	CustomDescriptionText               *string
	RestrictDirectMessage               *string
	RestrictTeamInvite                  *string
	RestrictPublicChannelManagement     *string
	RestrictPrivateChannelManagement    *string
	RestrictPublicChannelCreation       *string
	RestrictPrivateChannelCreation      *string
	RestrictPublicChannelDeletion       *string
	RestrictPrivateChannelDeletion      *string
	RestrictPrivateChannelManageMembers *string
	UserStatusAwayTimeout               *int64
	MaxChannelsPerTeam                  *int64
	MaxNotificationsPerChannel          *int64
}





type LocalizationSettings struct {
	DefaultServerLocale *string
	DefaultClientLocale *string
	AvailableLocales    *string
}





type Config struct {
	ServiceSettings      ServiceSettings
	SqlSettings          SqlSettings
	LogSettings          LogSettings
	PasswordSettings     PasswordSettings
	EmailSettings        EmailSettings
	RateLimitSettings    RateLimitSettings
	PrivacySettings      PrivacySettings
	GoogleSettings       SSOSettings
	LocalizationSettings LocalizationSettings
	ClusterSettings      ClusterSettings

}

func (o *Config) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (o *Config) GetSSOService(service string) *SSOSettings {
	switch service {
	case SERVICE_GOOGLE:
		return &o.GoogleSettings
	}

	return nil
}

func ConfigFromJson(data io.Reader) *Config {
	decoder := json.NewDecoder(data)
	var o Config
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

func (o *Config) SetDefaults() {

	if len(o.SqlSettings.AtRestEncryptKey) == 0 {
		o.SqlSettings.AtRestEncryptKey = NewRandomString(32)
	}



	if o.ServiceSettings.SiteURL == nil {
		o.ServiceSettings.SiteURL = new(string)
		*o.ServiceSettings.SiteURL = SERVICE_SETTINGS_DEFAULT_SITE_URL
	}


	if o.ServiceSettings.EnableDeveloper == nil {
		o.ServiceSettings.EnableDeveloper = new(bool)
		*o.ServiceSettings.EnableDeveloper = false
	}


	if o.PasswordSettings.MinimumLength == nil {
		o.PasswordSettings.MinimumLength = new(int)
		*o.PasswordSettings.MinimumLength = PASSWORD_MINIMUM_LENGTH
	}

	if o.PasswordSettings.Lowercase == nil {
		o.PasswordSettings.Lowercase = new(bool)
		*o.PasswordSettings.Lowercase = false
	}


	if o.EmailSettings.EnableSignInWithEmail == nil {
		o.EmailSettings.EnableSignInWithEmail = new(bool)

		if o.EmailSettings.EnableSignUpWithEmail == true {
			*o.EmailSettings.EnableSignInWithEmail = true
		} else {
			*o.EmailSettings.EnableSignInWithEmail = false
		}
	}

	if o.EmailSettings.EnableSignInWithUsername == nil {
		o.EmailSettings.EnableSignInWithUsername = new(bool)
		*o.EmailSettings.EnableSignInWithUsername = false
	}



	if o.ServiceSettings.SessionLengthWebInDays == nil {
		o.ServiceSettings.SessionLengthWebInDays = new(int)
		*o.ServiceSettings.SessionLengthWebInDays = 30
	}

	if o.ServiceSettings.SessionLengthMobileInDays == nil {
		o.ServiceSettings.SessionLengthMobileInDays = new(int)
		*o.ServiceSettings.SessionLengthMobileInDays = 30
	}

	if o.ServiceSettings.SessionLengthSSOInDays == nil {
		o.ServiceSettings.SessionLengthSSOInDays = new(int)
		*o.ServiceSettings.SessionLengthSSOInDays = 30
	}

	if o.ServiceSettings.SessionCacheInMinutes == nil {
		o.ServiceSettings.SessionCacheInMinutes = new(int)
		*o.ServiceSettings.SessionCacheInMinutes = 10
	}

	if o.ServiceSettings.EnableCommands == nil {
		o.ServiceSettings.EnableCommands = new(bool)
		*o.ServiceSettings.EnableCommands = false
	}

	if o.ServiceSettings.AllowCorsFrom == nil {
		o.ServiceSettings.AllowCorsFrom = new(string)
		*o.ServiceSettings.AllowCorsFrom = SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM
	}

	if o.ServiceSettings.WebserverMode == nil {
		o.ServiceSettings.WebserverMode = new(string)
		*o.ServiceSettings.WebserverMode = "gzip"
	} else if *o.ServiceSettings.WebserverMode == "regular" {
		*o.ServiceSettings.WebserverMode = "gzip"
	}

	if o.ServiceSettings.EnableCustomEmoji == nil {
		o.ServiceSettings.EnableCustomEmoji = new(bool)
		*o.ServiceSettings.EnableCustomEmoji = true
	}

	if o.ServiceSettings.RestrictCustomEmojiCreation == nil {
		o.ServiceSettings.RestrictCustomEmojiCreation = new(string)
		*o.ServiceSettings.RestrictCustomEmojiCreation = RESTRICT_EMOJI_CREATION_ALL
	}

	if o.ServiceSettings.RestrictPostDelete == nil {
		o.ServiceSettings.RestrictPostDelete = new(string)
		*o.ServiceSettings.RestrictPostDelete = PERMISSIONS_DELETE_POST_ALL
	}

	if o.ServiceSettings.AllowEditPost == nil {
		o.ServiceSettings.AllowEditPost = new(string)
		*o.ServiceSettings.AllowEditPost = ALLOW_EDIT_POST_ALWAYS
	}

	if o.ServiceSettings.PostEditTimeLimit == nil {
		o.ServiceSettings.PostEditTimeLimit = new(int)
		*o.ServiceSettings.PostEditTimeLimit = 300
	}

	if o.ClusterSettings.InterNodeListenAddress == nil {
		o.ClusterSettings.InterNodeListenAddress = new(string)
		*o.ClusterSettings.InterNodeListenAddress = ":8075"
	}

	if o.ClusterSettings.Enable == nil {
		o.ClusterSettings.Enable = new(bool)
		*o.ClusterSettings.Enable = false
	}

	if o.ClusterSettings.InterNodeUrls == nil {
		o.ClusterSettings.InterNodeUrls = []string{}
	}


	if o.LocalizationSettings.DefaultServerLocale == nil {
		o.LocalizationSettings.DefaultServerLocale = new(string)
		*o.LocalizationSettings.DefaultServerLocale = DEFAULT_LOCALE
	}

	if o.LocalizationSettings.DefaultClientLocale == nil {
		o.LocalizationSettings.DefaultClientLocale = new(string)
		*o.LocalizationSettings.DefaultClientLocale = DEFAULT_LOCALE
	}

	if o.LocalizationSettings.AvailableLocales == nil {
		o.LocalizationSettings.AvailableLocales = new(string)
		*o.LocalizationSettings.AvailableLocales = ""
	}

	if o.LogSettings.EnableDiagnostics == nil {
		o.LogSettings.EnableDiagnostics = new(bool)
		*o.LogSettings.EnableDiagnostics = true
	}

	if o.RateLimitSettings.Enable == nil {
		o.RateLimitSettings.Enable = new(bool)
		*o.RateLimitSettings.Enable = false
	}

	if o.RateLimitSettings.MaxBurst == nil {
		o.RateLimitSettings.MaxBurst = new(int)
		*o.RateLimitSettings.MaxBurst = 100
	}

	if o.ServiceSettings.ConnectionSecurity == nil {
		o.ServiceSettings.ConnectionSecurity = new(string)
		*o.ServiceSettings.ConnectionSecurity = ""
	}

	if o.ServiceSettings.TLSKeyFile == nil {
		o.ServiceSettings.TLSKeyFile = new(string)
		*o.ServiceSettings.TLSKeyFile = SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE
	}

	if o.ServiceSettings.TLSCertFile == nil {
		o.ServiceSettings.TLSCertFile = new(string)
		*o.ServiceSettings.TLSCertFile = SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE
	}

	if o.ServiceSettings.UseLetsEncrypt == nil {
		o.ServiceSettings.UseLetsEncrypt = new(bool)
		*o.ServiceSettings.UseLetsEncrypt = false
	}

	if o.ServiceSettings.LetsEncryptCertificateCacheFile == nil {
		o.ServiceSettings.LetsEncryptCertificateCacheFile = new(string)
		*o.ServiceSettings.LetsEncryptCertificateCacheFile = "./config/letsencrypt.cache"
	}

	if o.ServiceSettings.ReadTimeout == nil {
		o.ServiceSettings.ReadTimeout = new(int)
		*o.ServiceSettings.ReadTimeout = SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT
	}

	if o.ServiceSettings.WriteTimeout == nil {
		o.ServiceSettings.WriteTimeout = new(int)
		*o.ServiceSettings.WriteTimeout = SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT
	}

	if o.ServiceSettings.Forward80To443 == nil {
		o.ServiceSettings.Forward80To443 = new(bool)
		*o.ServiceSettings.Forward80To443 = false
	}


	if o.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds == nil {
		o.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds = new(int64)
		*o.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds = 5000
	}

	if o.ServiceSettings.ClusterLogTimeoutMilliseconds == nil {
		o.ServiceSettings.ClusterLogTimeoutMilliseconds = new(int)
		*o.ServiceSettings.ClusterLogTimeoutMilliseconds = 2000
	}

}

func (o *Config) IsValid() *AppError {

	if o.ServiceSettings.MaximumLoginAttempts <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.login_attempts.app_error", nil, "")
	}

	if len(*o.ServiceSettings.SiteURL) != 0 {
		if _, err := url.ParseRequestURI(*o.ServiceSettings.SiteURL); err != nil {
			return NewLocAppError("Config.IsValid", "model.config.is_valid.site_url.app_error", nil, "")
		}
	}

	if len(o.ServiceSettings.ListenAddress) == 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.listen_address.app_error", nil, "")
	}

	if *o.ClusterSettings.Enable && *o.EmailSettings.EnableEmailBatching {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.cluster_email_batching.app_error", nil, "")
	}

	if len(*o.ServiceSettings.SiteURL) == 0 && *o.EmailSettings.EnableEmailBatching {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.site_url_email_batching.app_error", nil, "")
	}


	if len(o.SqlSettings.AtRestEncryptKey) < 32 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.encrypt_sql.app_error", nil, "")
	}

	if !(o.SqlSettings.DriverName == DATABASE_DRIVER_MYSQL || o.SqlSettings.DriverName == DATABASE_DRIVER_POSTGRES) {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.sql_driver.app_error", nil, "")
	}

	if o.SqlSettings.MaxIdleConns <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.sql_idle.app_error", nil, "")
	}

	if len(o.SqlSettings.DataSource) == 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.sql_data_src.app_error", nil, "")
	}

	if o.SqlSettings.MaxOpenConns <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.sql_max_conn.app_error", nil, "")
	}


	if o.RateLimitSettings.MemoryStoreSize <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.rate_mem.app_error", nil, "")
	}

	if o.RateLimitSettings.PerSec <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.rate_sec.app_error", nil, "")
	}

	if *o.PasswordSettings.MinimumLength < PASSWORD_MINIMUM_LENGTH || *o.PasswordSettings.MinimumLength > PASSWORD_MAXIMUM_LENGTH {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.password_length.app_error", map[string]interface{}{"MinLength": PASSWORD_MINIMUM_LENGTH, "MaxLength": PASSWORD_MAXIMUM_LENGTH}, "")
	}



	if *o.RateLimitSettings.MaxBurst <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.max_burst.app_error", nil, "")
	}


	if !(*o.ServiceSettings.ConnectionSecurity == CONN_SECURITY_NONE || *o.ServiceSettings.ConnectionSecurity == CONN_SECURITY_TLS) {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.webserver_security.app_error", nil, "")
	}

	if *o.ServiceSettings.ReadTimeout <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.read_timeout.app_error", nil, "")
	}

	if *o.ServiceSettings.WriteTimeout <= 0 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.write_timeout.app_error", nil, "")
	}

	if *o.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds < 1000 {
		return NewLocAppError("Config.IsValid", "model.config.is_valid.time_between_user_typing.app_error", nil, "")
	}

	return nil
}

func (o *Config) GetSanitizeOptions() map[string]bool {
	options := map[string]bool{}
	options["fullname"] = o.PrivacySettings.ShowFullName
	options["email"] = o.PrivacySettings.ShowEmailAddress

	return options
}

func (o *Config) Sanitize() {

	o.EmailSettings.InviteSalt = FAKE_SETTING
	o.EmailSettings.PasswordResetSalt = FAKE_SETTING
	if len(o.EmailSettings.SMTPPassword) > 0 {
		o.EmailSettings.SMTPPassword = FAKE_SETTING
	}


	o.SqlSettings.DataSource = FAKE_SETTING
	o.SqlSettings.AtRestEncryptKey = FAKE_SETTING

	for i := range o.SqlSettings.DataSourceReplicas {
		o.SqlSettings.DataSourceReplicas[i] = FAKE_SETTING
	}

	for i := range o.SqlSettings.DataSourceSearchReplicas {
		o.SqlSettings.DataSourceSearchReplicas[i] = FAKE_SETTING
	}
}
