package domain

import "errors"

// Plan An assessment plan, such as those provided by a FedRAMP assessor.
// Here are some real-world examples for Assets, Platforms, Subjects and Inventory Items within an OSCAL Assessment Plan:
// 1. Assets: This could be something like a customer database within a retail company. It's an asset because it's crucial to the business operation, storing all the essential customer details such as addresses, contact information, and purchase history.
// 2. Platforms: This could be the retail company's online E-commerce platform which hosts their online store, and where transactions occur. The platform might involve web servers, database servers, or a cloud environment.
// 3. Subjects: If the company is performing a security assessment, the subject could be the encryption method or security protocols used to protect the customer data while in transit or at rest in the database.
// 4. Inventory Items: These could be the individual servers or workstations used within the company. Inventory workstations are the physical machines or software applications used by employees that may have vulnerabilities or exposure to risk that need to be tracked and mitigated.
//
// Relation between Tasks, Activities and Steps:
//
// Scenario: Conducting a cybersecurity assessment of an organization's systems.
//
// 1. Task: The major task could be "Conduct vulnerability scanning on servers."
// 2. Activity: Within this task, an activity could be "Prepare servers for vulnerability scan."
// 3. Step: The steps that make up this activity could be things like:
//   - "Identify all servers"
//   - "Ensure necessary permissions are in place for scanning"
//   - "Check that scanning software is properly installed and updated."
//
// Another activity under the same task could be "Execute vulnerability scanning," and steps for that activity might include:
//
// 1. "Begin scanning process through scanning software."
// 2. "Monitor progress of scan."
// 3. "Document any issues or vulnerabilities identified."
//
// The process would continue like this with tasks broken down into activities, and activities broken down into steps.
//
// These concepts still apply in the context of automated tools or systems. In fact, the OSCAL model is designed to support both manual and automated processes.
// 1.	Task: The major task could be “Automated Compliance Checking”
// 2.	Activity: This task could have multiple activities such as:
// ▪	“Configure Automated Tool with necessary parameters”
// ▪	“Run Compliance Check”
// ▪	“Collect and Analyze Compliance Data”
// 3.	Step: In each of these activities, there are several subprocesses or actions (Steps). For example, under “Configure Automated Tool with necessary parameters”, the steps could be:
// ▪	“Define the criteria based on selected standards”
// ▪	“Set the scope or target systems for the assessment”
// ▪	“Specify the output (report) format”
// In context of an automated compliance check, the description of Task, Activity, and Step provides a systematic plan or procedure that the tool is expected to follow. This breakdown of tasks, activities, and steps could also supply useful context and explain the tool’s operation and results to system admins, auditors or other stakeholders. It also allows for easier troubleshooting in the event of problems.
type Plan struct {
	Uuid Uuid `json:"uuid"`

	// Title A name given to the assessment plan. OSCAL doesn't have this, but we need it for our use case.
	Title string `json:"title,omitempty"`

	// We might switch to struct embedding for fields like Metadata, Props, etc.
	Metadata Metadata `json:"metadata"`

	// Assets Identifies the assets used to perform this assessment, such as the assessment team, scanning tools, and assumptions. Mostly CF in our case.
	Assets Assets `json:"assets"`

	// BackMatter A collection of resources that may be referenced from within the OSCAL document instance.
	BackMatter BackMatter `json:"backMatter"`

	// Reference to a System Security Plan
	ImportSSP Uuid `json:"importSSP"`

	// LocalDefinitions Used to define data objects that are used in the assessment plan, that do not appear in the referenced SSP.
	// Reference to LocalDefinition
	LocalDefinitions LocalDefinition `json:"localDefinitions"`

	// ReviewedControls Identifies the controls being assessed and their control objectives.
	ReviewedControls []ControlsAndObjectives `json:"reviewedControls"`

	// Tasks Represents a scheduled event or milestone, which may be associated with a series of assessment actions.
	Tasks []Task `json:"tasks"`

	// TermsAndConditions Used to define various terms and conditions under which an assessment, described by the plan, can be performed. Each child part defines a different type of term or condition.
	TermsAndConditions []Part `json:"termsAndConditions"`
}

func NewPlan() *Plan {
	revision := NewRevision("Initial version", "Initial version", "")

	metadata := Metadata{
		Revisions: []Revision{revision},
		Actions: []Action{
			{
				Uuid:  NewUuid(),
				Title: "Create",
			},
		},
	}

	return &Plan{
		Uuid:     NewUuid(),
		Metadata: metadata,
		Assets: Assets{
			Components: []Uuid{},
			Platforms:  []Uuid{},
		},
	}
}

func (p *Plan) AddAsset(assetUuid Uuid, assetType string) {
	if assetType == "component" {
		p.Assets.Components = append(p.Assets.Components, assetUuid)
	} else if assetType == "platform" {
		p.Assets.Platforms = append(p.Assets.Components, assetUuid)
	}
}

func (p *Plan) GetTask(id Uuid) *Task {
	for _, task := range p.Tasks {
		if task.Uuid == id {
			return &task
		}
	}
	return nil
}

func (p *Plan) CreateTask(task Task) error {
	// Validate the task
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}

	if task.Type != TaskTypeMilestone && task.Type != TaskTypeAction {
		return errors.New("task type must be either 'milestone' or 'action'")
	}

	// Add the task to the Tasks slice
	p.Tasks = append(p.Tasks, task)

	return nil
}

func (p *Plan) Ready() bool {
	// If there are no Tasks then there's nothing to run.
	if len(p.Tasks) == 0 {
		return false
	}

	// Check if the tasks have a schedule for the future and at least one subject.
	for _, task := range p.Tasks {
		if task.Subjects.Query == "" && len(task.Subjects.Labels) == 0 && len(task.Subjects.Ids) == 0 && len(task.Subjects.Expressions) == 0 && len(task.Schedule) == 0 {
			continue
		}
	}

	return false
}

func (p *Plan) JobSpecification() (JobSpecification, error) {
	// Step 1: Create an empty JobSpecification
	jobSpec := JobSpecification{}

	// Step 2: Set RuntimeUuid to a new UUID
	jobSpec.RuntimeUuid = NewUuid().String()

	// Step 3: Create an AssessmentPlanInformation from the Plan
	planInfo := AssessmentPlanInformation{
		Uuid:  p.Uuid.String(),
		Title: p.Title,
	}

	// Step 3.2: For each Task in the Plan, create a TaskInformation
	for _, task := range p.Tasks {
		taskInfo := TaskInformation{
			Uuid:     task.Uuid.String(),
			Title:    task.Title,
			Selector: task.Subjects,
		}

		// Step 3.2.2: For each Activity in the Task, create an ActivityInformation
		for _, activity := range task.Activities {
			activityInfo := ActivityInformation{
				Uuid:     activity.Uuid.String(),
				Title:    activity.Title,
				Provider: activity.Provider,
			}
			taskInfo.Activities = append(taskInfo.Activities, activityInfo)
		}

		planInfo.Tasks = append(planInfo.Tasks, taskInfo)
	}

	// Step 4: Set Plan in the JobSpecification to the created AssessmentPlanInformation
	jobSpec.Plan = planInfo

	// Step 5: Return the JobSpecification
	return jobSpec, nil
}

type TaskType string

const (
	TaskTypeMilestone TaskType = "milestone"
	TaskTypeAction    TaskType = "action"
)

type Task struct {
	Uuid             Uuid             `json:"uuid"`
	Title            string           `json:"title,omitempty"`
	Description      string           `json:"description,omitempty"`
	Props            []Property       `json:"props,omitempty"`
	Links            []Link           `json:"links,omitempty"`
	Remarks          string           `json:"remarks,omitempty"`
	Type             TaskType         `json:"type"`
	Activities       []Activity       `json:"activities"`
	Dependencies     []TaskDependency `json:"dependencies"`
	ResponsibleRoles []Uuid           `json:"responsibleRoles"`
	Subjects         SubjectSelection `json:"subjects"`
	Tasks            []Uuid           `json:"tasks"`
	Schedule         []string         `json:"schedule"`
}

func (t *Task) AddSubject(subjects SubjectSelection) error {
	// Validate the subjects
	if subjects.Title == "" {
		return errors.New("subjects title cannot be empty")
	}

	// Check if only one of Query, Labels, Expressions, and Ids is set
	fieldsSet := 0
	if len(subjects.Ids) > 0 {
		fieldsSet++
	}
	if subjects.Query != "" {
		fieldsSet++
	}
	if len(subjects.Expressions) > 0 {
		fieldsSet++
	}
	if len(subjects.Labels) > 0 {
		fieldsSet++
	}

	// If more than one is set, unset the others based on the priority order
	if fieldsSet > 1 {
		if len(subjects.Ids) > 0 {
			subjects.Query = ""
			subjects.Expressions = nil
			subjects.Labels = nil
		} else if subjects.Query != "" {
			subjects.Expressions = nil
			subjects.Labels = nil
		} else if len(subjects.Expressions) > 0 {
			subjects.Labels = nil
		}
	}

	// Add the subjects to the Subjects slice
	t.Subjects = subjects

	return nil
}

func (t *Task) AddActivity(activity Activity) error {
	// Validate the activity
	if activity.Title == "" {
		return errors.New("activity title cannot be empty")
	}

	// Add the activity to the Activities slice
	t.Activities = append(t.Activities, activity)

	return nil
}

type TaskDependency struct {
	TaskId  Uuid   `json:"taskUuid"`
	Remarks string `json:"remarks"`
}

// Assets Identifies the assets used to perform this assessment, such as the assessment team, scanning tools, and assumptions.
type Assets struct {
	// Reference to component.Component
	Components []Uuid `json:"components"`

	// Used to represent the toolset used to perform aspects of the assessment.
	Platforms []Uuid `json:"platforms"`
}

type Platform struct {
	Uuid        Uuid       `json:"uuid"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Props       []Property `json:"props,omitempty"`

	Links   []Link `json:"links,omitempty"`
	Remarks string `json:"remarks,omitempty"`

	// Reference to component.Component
	UsesComponents []Uuid `json:"usesComponents"`
}

type ControlsAndObjectives struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Props       []Property `json:"props,omitempty"`

	Links   []Link `json:"links,omitempty"`
	Remarks string `json:"remarks,omitempty"`

	Objectives        []ObjectiveSelection `json:"objectives"`
	ControlSelections Selection            `json:"controlSelections"`
}

type ObjectiveSelection struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Props       []Property `json:"props,omitempty"`

	Links      []Link `json:"links,omitempty"`
	Remarks    string `json:"remarks,omitempty"`
	IncludeAll bool   `json:"includeAll"`
	Exclude    []Uuid `json:"exclude"`
	Include    []Uuid `json:"include"`
}

type LocalDefinition struct {
	Remarks string `json:"remarks,omitempty"`

	// Reference to Activity
	Activities []Uuid `json:"activities"`

	// Reference to component.Component
	Components []Uuid `json:"components"`

	// Reference to ssp.InventoryItem
	InventoryItems []Uuid `json:"inventoryItems"`

	Objectives []Objective `json:"objectives"`

	// Reference to identity.User
	Users []Uuid `json:"users"`
}

// Objective A local objective is a security control or requirement that is specific to the system or organization under assessment.
type Objective struct {
	Uuid        Uuid       `json:"uuid"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Props       []Property `json:"props,omitempty"`

	Links   []Link `json:"links,omitempty"`
	Remarks string `json:"remarks,omitempty"`
	Parts   []Part `json:"parts,omitempty"`

	Control Uuid `json:"control"`
}

type SubjectType string

const (
	SubjectTypeComponent     SubjectType = "component"
	SubjectTypeInventoryItem SubjectType = "inventoryItem"
	SubjectTypeLocation      SubjectType = "location"
	SubjectTypeParty         SubjectType = "party"
	SubjectTypeUser          SubjectType = "user"
)

// Subject Identifies system elements being assessed, such as components, inventory items, and locations.
// In the assessment plan, this identifies a planned assessment subject.
// In the assessment results this is an actual assessment subject, and reflects any changes from the plan. exactly what will be the focus of this assessment.
type Subject struct {
	Uuid        Uuid        `json:"uuid"`
	Type        SubjectType `json:"type"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Props       []Property  `json:"props,omitempty"`
	Links       []Link      `json:"links,omitempty"`
	Remarks     string      `json:"remarks,omitempty"`
}

// SubjectSelection Identifies system elements being assessed, such as components, inventory items, and locations by specifying a selection criteria.
// We do not directly store SubjectIds as we might not know the actual subjects before running the assessment.
// The assessment runtime evaluates the selection by running the providers and returns back with subject ids.
type SubjectSelection struct {
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Query       string                   `yaml:"query" json:"query"`
	Labels      map[string]string        `yaml:"labels,omitempty" json:"labels,omitempty"`
	Expressions []SubjectMatchExpression `yaml:"expressions,omitempty" json:"expressions,omitempty"`
	Ids         []string                 `yaml:"ids,omitempty" json:"ids,omitempty"`
}

type SubjectMatchExpression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type Activity struct {
	Uuid             Uuid                  `json:"uuid"`
	Title            string                `json:"title,omitempty"`
	Description      string                `json:"description,omitempty"`
	Props            []Property            `json:"props,omitempty"`
	Links            []Link                `json:"links,omitempty"`
	Remarks          string                `json:"remarks,omitempty"`
	ResponsibleRoles []Uuid                `json:"responsibleRoles"`
	Steps            []Step                `json:"steps"`
	Provider         ProviderConfiguration `json:"provider"`
}

type Step struct {
	Uuid             Uuid        `json:"uuid"`
	Title            string      `json:"title,omitempty"`
	Description      string      `json:"description,omitempty"`
	Props            []Property  `json:"props,omitempty"`
	Links            []Link      `json:"links,omitempty"`
	Remarks          string      `json:"remarks,omitempty"`
	ResponsibleRoles []Uuid      `json:"responsibleRoles"`
	Objectives       []Objective `json:"objectives"`
}
