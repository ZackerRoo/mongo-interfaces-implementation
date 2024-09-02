package model

type Knowledge struct {
	ID                  string   `bson:"_id,omitempty" json:"id"`
	Title               string   `bson:"title,omitempty" json:"title"`
	Tags                []string `bson:"tags,omitempty" json:"tags"`
	TechniquesID        []string `bson:"techniquesId,omitempty" json:"techniquesId"`
	TacticsID           []string `bson:"tacticsId,omitempty" json:"tacticsId"`
	KnowledgeType       []string `bson:"knowledgeType,omitempty" json:"knowledgeType"`
	KnowledgeSource     []string `bson:"knowledgeSource,omitempty" json:"knowledgeSource"`
	Confidentiality     string   `bson:"confidentiality,omitempty" json:"confidentiality"`
	Abstract            string   `bson:"abstract,omitempty" json:"abstract"`
	Content             string   `bson:"content,omitempty" json:"content"`
	Detection           string   `bson:"detection,omitempty" json:"detection"`
	Mitigations         string   `bson:"mitigations,omitempty" json:"mitigations"`
	Recommendations     string   `bson:"recommendations,omitempty" json:"recommendations"`
	Directory           string   `bson:"directory,omitempty" json:"directory"`
	Techniques          string   `bson:"techniques,omitempty" json:"techniques"`
	Tactics             string   `bson:"tactics,omitempty" json:"tactics"`
	UsedExploits        string   `bson:"usedExploits,omitempty" json:"usedExploits"`
	Alias               string   `bson:"alias,omitempty" json:"alias"`
	VulType             string   `bson:"vulType,omitempty" json:"vulType"`
	Affiliation         string   `bson:"affiliation,omitempty" json:"affiliation"`
	UsedTools           string   `bson:"usedTools,omitempty" json:"usedTools"`
	StrategicCapability string   `bson:"strategicCapability,omitempty" json:"strategicCapability"`
	FirstActivity       string   `bson:"firstActivity,omitempty" json:"firstActivity"`
	LatestActivity      string   `bson:"latestActivity,omitempty" json:"latestActivity"`
	TargetedGeography   string   `bson:"targetedGeography,omitempty" json:"targetedGeography"`
	TimeLine            string   `bson:"timeLine,omitempty" json:"timeLine"`
	Scenario            string   `bson:"scenario,omitempty" json:"scenario"`
	Motivations         string   `bson:"motivations,omitempty" json:"motivations"`
	TargetedIndustry    string   `bson:"targetedIndustry,omitempty" json:"targetedIndustry"`
	Preparation         string   `bson:"preparation,omitempty" json:"preparation"`
	Alert               string   `bson:"alert,omitempty" json:"alert"`
	Analysis            string   `bson:"analysis,omitempty" json:"analysis"`
	Traces              string   `bson:"traces,omitempty" json:"traces"`
	Containment         string   `bson:"containment,omitempty" json:"containment"`
	Eradication         string   `bson:"eradication,omitempty" json:"eradication"`
	Recovery            string   `bson:"recovery,omitempty" json:"recovery"`
	FollowUp            string   `bson:"followUp,omitempty" json:"followUp"`
	DisposalProcess     string   `bson:"disposalProcess,omitempty" json:"disposalProcess"`
	Cases               string   `bson:"cases,omitempty" json:"cases"`
	Cve                 string   `bson:"cve,omitempty" json:"cve"`
	Cnnvd               string   `bson:"cnnvd,omitempty" json:"cnnvd"`
	Cwd                 string   `bson:"cwd,omitempty" json:"cwd"`
	Cvss                string   `bson:"cvss,omitempty" json:"cvss"`
	Bugtraq             string   `bson:"bugtraq,omitempty" json:"bugtraq"`
	CvssStr             string   `bson:"cvssStr,omitempty" json:"cvssStr"`
	Msf                 string   `bson:"msf,omitempty" json:"msf"`
	Exploitdb           string   `bson:"exploitdb,omitempty" json:"exploitdb"`
	IsExp               string   `bson:"isExp,omitempty" json:"isExp"`
	Vendor              string   `bson:"vendor,omitempty" json:"vendor"`
	AppType             string   `bson:"appType,omitempty" json:"appType"`
	Consequence         string   `bson:"consequence,omitempty" json:"consequence"`
	FingerPrint         string   `bson:"fingerPrint,omitempty" json:"fingerPrint"`
	RevisionDate        []string `bson:"revisionDate,omitempty" json:"revisionDate"`
	Products            string   `bson:"products,omitempty" json:"products"`
	Reference           string   `bson:"reference,omitempty" json:"reference"`
	Author              []string `bson:"author,omitempty" json:"author"`
	UID                 string   `bson:"uid,omitempty" json:"uid"`
	SubTechniquesID     []string `bson:"subTechniquesId,omitempty" json:"subTechniquesId"`
	Platforms           []string `bson:"platforms,omitempty" json:"platforms"`
	AffectedVerison     string   `bson:"affectedVerison,omitempty" json:"affectedVerison"`
	ThreatSeverity      string   `bson:"threatSeverity,omitempty" json:"threatSeverity"`
	Solution            string   `bson:"solution,omitempty" json:"solution"`
	Cnvd                string   `bson:"cnvd,omitempty" json:"cnvd"`
	Cwe                 string   `bson:"cwe,omitempty" json:"cwe"`
	AppName             string   `bson:"appName,omitempty" json:"appName"`
	OrganizationIds     string   `bson:"organizationIds,omitempty" json:"organizationIds"`
	IoC                 string   `bson:"ioc,omitempty" json:"IoC"`
	TiName              string   `bson:"tiName,omitempty" json:"TiName"`
	InputParameters     string   `bson:"inputParameters,omitempty" json:"inputParameters"`
	OutputParameters    string   `bson:"outputParameters,omitempty" json:"outputParameters"`
	Success             bool     `bson:"success,omitempty" json:"success"`
	Message             string   `bson:"message,omitempty" json:"message"`
}

type KnowledgeFilter struct {
	Abstract        string             `bson:"abstract,omitempty"`
	Confidentiality string             `bson:"confidentiality,omitempty"`
	Consequence     string             `bson:"consequence,omitempty"`
	Content         string             `bson:"content,omitempty"`
	Detection       string             `bson:"detection,omitempty"`
	ID              string             `bson:"id,omitempty"`
	KnowledgeSource []string           `bson:"knowledgeSource,omitempty"`
	KnowledgeType   []string           `bson:"knowledgeType,omitempty"`
	Mitigations     string             `bson:"mitigations,omitempty"`
	Name            string             `bson:"name,omitempty"`
	Recommendations string             `bson:"recommendations,omitempty"`
	Solution        string             `bson:"solution,omitempty"`
	TacticsID       []string           `bson:"tacticsId,omitempty"`
	Tags            []string           `bson:"tags,omitempty"`
	TechniquesID    []string           `bson:"techniquesId,omitempty"`
	Title           string             `bson:"title,omitempty"`
	Tactics         string             `bson:"tactics,omitempty"`
	SubTechniquesID []string           `bson:"subTechniquesId,omitempty"`
	AND             []*KnowledgeFilter `bson:"AND,omitempty"`
}

type NewKnowledge struct {
	ID                  string   `bson:"_id,omitempty" json:"id"`
	Title               string   `bson:"title,omitempty" json:"title"`
	Tags                []string `bson:"tags,omitempty" json:"tags"`
	TechniquesID        []string `bson:"techniquesId,omitempty" json:"techniquesId"`
	TacticsID           []string `bson:"tacticsId,omitempty" json:"tacticsId"`
	KnowledgeType       []string `bson:"knowledgeType,omitempty" json:"knowledgeType"`
	KnowledgeSource     []string `bson:"knowledgeSource,omitempty" json:"knowledgeSource"`
	Confidentiality     string   `bson:"confidentiality,omitempty" json:"confidentiality"`
	Abstract            string   `bson:"abstract,omitempty" json:"abstract"`
	Content             string   `bson:"content,omitempty" json:"content"`
	Detection           string   `bson:"detection,omitempty" json:"detection"`
	Mitigations         string   `bson:"mitigations,omitempty" json:"mitigations"`
	Recommendations     string   `bson:"recommendations,omitempty" json:"recommendations"`
	Directory           string   `bson:"directory,omitempty" json:"directory"`
	Techniques          string   `bson:"techniques,omitempty" json:"techniques"`
	Tactics             string   `bson:"tactics,omitempty" json:"tactics"`
	UsedExploits        string   `bson:"usedExploits,omitempty" json:"usedExploits"`
	Alias               string   `bson:"alias,omitempty" json:"alias"`
	VulType             string   `bson:"vulType,omitempty" json:"vulType"`
	Affiliation         string   `bson:"affiliation,omitempty" json:"affiliation"`
	UsedTools           string   `bson:"usedTools,omitempty" json:"usedTools"`
	StrategicCapability string   `bson:"strategicCapability,omitempty" json:"strategicCapability"`
	FirstActivity       string   `bson:"firstActivity,omitempty" json:"firstActivity"`
	LatestActivity      string   `bson:"latestActivity,omitempty" json:"latestActivity"`
	TargetedGeography   string   `bson:"targetedGeography,omitempty" json:"targetedGeography"`
	TimeLine            string   `bson:"timeLine,omitempty" json:"timeLine"`
	Scenario            string   `bson:"scenario,omitempty" json:"scenario"`
	Motivations         string   `bson:"motivations,omitempty" json:"motivations"`
	TargetedIndustry    string   `bson:"targetedIndustry,omitempty" json:"targetedIndustry"`
	Preparation         string   `bson:"preparation,omitempty" json:"preparation"`
	Alert               string   `bson:"alert,omitempty" json:"alert"`
	Analysis            string   `bson:"analysis,omitempty" json:"analysis"`
	Traces              string   `bson:"traces,omitempty" json:"traces"`
	Containment         string   `bson:"containment,omitempty" json:"containment"`
	Eradication         string   `bson:"eradication,omitempty" json:"eradication"`
	Recovery            string   `bson:"recovery,omitempty" json:"recovery"`
	FollowUp            string   `bson:"followUp,omitempty" json:"followUp"`
	DisposalProcess     string   `bson:"disposalProcess,omitempty" json:"disposalProcess"`
	Cases               string   `bson:"cases,omitempty" json:"cases"`
	Cve                 string   `bson:"cve,omitempty" json:"cve"`
	Cnnvd               string   `bson:"cnnvd,omitempty" json:"cnnvd"`
	Cwd                 string   `bson:"cwd,omitempty" json:"cwd"`
	Cvss                string   `bson:"cvss,omitempty" json:"cvss"`
	Bugtraq             string   `bson:"bugtraq,omitempty" json:"bugtraq"`
	CvssStr             string   `bson:"cvssStr,omitempty" json:"cvssStr"`
	Msf                 string   `bson:"msf,omitempty" json:"msf"`
	Exploitdb           string   `bson:"exploitdb,omitempty" json:"exploitdb"`
	IsExp               string   `bson:"isExp,omitempty" json:"isExp"`
	Vendor              string   `bson:"vendor,omitempty" json:"vendor"`
	AppType             string   `bson:"appType,omitempty" json:"appType"`
	Consequence         string   `bson:"consequence,omitempty" json:"consequence"`
	FingerPrint         string   `bson:"fingerPrint,omitempty" json:"fingerPrint"`
	RevisionDate        []string `bson:"revisionDate,omitempty" json:"revisionDate"`
	Products            string   `bson:"products,omitempty" json:"products"`
	Reference           string   `bson:"reference,omitempty" json:"reference"`
	Author              []string `bson:"author,omitempty" json:"author"`
	UID                 string   `bson:"uid,omitempty" json:"uid"`
	SubTechniquesID     []string `bson:"subTechniquesId,omitempty" json:"subTechniquesId"`
	Platforms           []string `bson:"platforms,omitempty" json:"platforms"`
	AffectedVerison     string   `bson:"affectedVerison,omitempty" json:"affectedVerison"`
	ThreatSeverity      string   `bson:"threatSeverity,omitempty" json:"threatSeverity"`
	Solution            string   `bson:"solution,omitempty" json:"solution"`
	Cnvd                string   `bson:"cnvd,omitempty" json:"cnvd"`
	Cwe                 string   `bson:"cwe,omitempty" json:"cwe"`
	AppName             string   `bson:"appName,omitempty" json:"appName"`
	OrganizationIds     string   `bson:"organizationIds,omitempty" json:"organizationIds"`
	IoC                 string   `bson:"ioc,omitempty" json:"IoC"`
	TiName              string   `bson:"tiName,omitempty" json:"TiName"`
	InputParameters     string   `bson:"inputParameters,omitempty" json:"inputParameters"`
	OutputParameters    string   `bson:"outputParameters,omitempty" json:"outputParameters"`
}

type DeletionStatus struct {
	Success bool   `bson:"success,omitempty" json:"success"`
	Message string `bson:"message,omitempty" json:"message"`
}
