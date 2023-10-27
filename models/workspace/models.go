package workspace

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	WorkspaceOrg      = "workspaceOrg"
	WorkSpacePersonal = "workspacePersonal"
)

type Workspace struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	ClerkOrgId  string             `json:"clerkOrgId" bson:"clerkOrgId"`
	Extensions  []interface{}      `json:"extensions" bson:"extensions"`
	Connections []string           `json:"connections" bson:"connections,omitempty"`
}

type WorkspaceMember struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	WorkspaceId string             `json:"workspaceId" bson:"workspaceId"`
	UserId      string             `json:"userId" bson:"userId"`
	Roles       []string           `json:"roles" bson:"roles"`
	Permissions []string           `json:"permissions" bson:"permissions"`
}

type Workspaces []Workspace

const (
	Teams    string = "teams"
	Stats    string = "stats"
	Insights string = "insights"
	Calendar string = "calendar"
	Memo     string = "memo"
	Planner  string = "planner"
)

const (
	GoogleCalendar string = "google_calendar"
	MsCalendar     string = "ms_calendar"
	Slack          string = "slack"
	Jira           string = "jira"
	ClickUp        string = "clickup"
	TimeDoctor     string = "time_doctor"
)

const (
	// Read resources
	TeamsRead    string = "teams:read"
	StatsRead    string = "stats:read"
	InsightsRead string = "insights:read"
	CalendarRead string = "calendar:read"
	MemoRead     string = "memo:read"
	PlannerRead  string = "planner:read"

	// Write resources
	TeamsWrite    string = "teams:write"
	CalendarWrite string = "calendar:write"
	MemoWrite     string = "memo:write"
	PlannerWrite  string = "planner:write"
)

var PermissionsAdminOrg = []string{
	TeamsRead,
	TeamsWrite,
	StatsRead,
	InsightsRead,
	CalendarRead,
	CalendarWrite,
	MemoRead,
	MemoWrite,
	PlannerRead,
	PlannerWrite,
}

const (
	WorkspaceRoleAdmin string = "admin"
)

type WorkspaceResponse struct {
	Workspace          Workspace `json:"workspaces"`
	Connections        []string  `json:"connections"`
	RoleBasedResources []string  `json:"roleBasedResources"`
}
