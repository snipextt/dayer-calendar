package workspace

import (
	"github.com/snipextt/dayer/models/connection"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	WorkspaceOrg      = "workspaceOrg"
	WorkSpacePersonal = "workspacePersonal"
)

type Workspace struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DefaultTeam primitive.ObjectID `json:"defaultTeam" bson:"defaultTeam"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	ClerkOrgId  string             `json:"clerkOrgId" bson:"clerkOrgId"`
	Extensions  []string           `json:"extensions" bson:"extensions"`
}

type MemberMeta struct {
	Source       string `json:"source" bson:"source"`
	TimeDoctorId string `json:"timeDoctorId" bson:"timeDoctorId"`
}

type Member struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Image       string             `json:"image" bson:"image"`
	Workspace   any                `json:"workspace" bson:"workspace"`
	User        any                `json:"user" bson:"user"`
	Manager     any                `json:"manager" bson:"manager"`
	Teams       any                `json:"teams" bson:"teams"`
	Roles       []string           `json:"roles" bson:"roles"`
	Permissions []string           `json:"permissions" bson:"permissions"`
	Meta        MemberMeta         `json:"meta" bson:"meta"`
}

type Team struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Workspace   any                `json:"workspace" bson:"workspace"`
	Owner       primitive.ObjectID `json:"owner" bson:"owner"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
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
	TeamsWrite      string = "teams:write"
	CalendarWrite   string = "calendar:write"
	MemoWrite       string = "memo:write"
	PlannerWrite    string = "planner:write"
	ManageWorkspace string = "workspace:write"
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
	ManageWorkspace,
	PlannerWrite,
}

const (
	WorkspaceRoleAdmin  string = "admin"
	WorkspaceRoleMember string = "member"
)

type WorkspaceEvent struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type WorkspaceResponse struct {
	Id                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PendingConnections []WorkspaceEvent   `json:"pendingConnections"`
	PendingActions     []WorkspaceEvent   `json:"pendingActions"`
	Teams              []Team             `json:"teams" bson:"teams"`
	RoleBasedResources []string           `json:"roleBasedResources"`
}

type WorkspaceAggregation struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	User        Member             `json:"user" bson:"user"`
	Extensions  []string           `json:"extensions" bson:"extensions"`
	Connections []connection.Model `json:"connections" bson:"connections"`
	Teams       []Team             `json:"teams" bson:"teams"`
}
