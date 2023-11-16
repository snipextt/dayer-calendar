package timedoctor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/snipextt/dayer/models"
	"github.com/snipextt/dayer/utils"
)

type TimeDoctorImageData = models.TimeDoctorImageData
type TimeDoctorActivity = models.TimeDoctorActivity

type TimeDoctorPagination struct {
	Current    string `json:"cur"`
	Next       string `json:"next"`
	Limit      int    `json:"limit"`
	TotalCount int    `json:"totalCount"`
}

type TimeDoctorUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TimeDoctorCompany struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	CompanyTimezone string `json:"companyTimezone"`
	UserCount       int    `json:"userCount"`
}

type TimeDoctorLoginAuthorizationResponseData struct {
	Companies []TimeDoctorCompany `json:"companies"`
	Email     string              `json:"email"`
	Token     string              `json:"token"`
	ExpiresAt string              `json:"expiresAt"`
}

type TimeDoctorResponse[T any] struct {
	Data    T                    `json:"data"`
	Error   string               `json:"error"`
	Message string               `json:"message"`
	Page    TimeDoctorPagination `json:"page"`
}

const baseurl1_0 = "https://api2.timedoctor.com/api/1.0"
const baseurl1_1 = "https://api2.timedoctor.com/api/1.1"

func Login(email, password, totp string) (response *TimeDoctorResponse[TimeDoctorLoginAuthorizationResponseData], err error) {
	url := baseurl1_0 + "/authorization/login"
	body := map[string]string{
		"email":       email,
		"password":    password,
		"totpCode":    totp,
		"permissions": "read",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	rsp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	err = json.NewDecoder(rsp.Body).Decode(&response)
	return
}

func GetCompanyUsers(token, companyId string) (response *TimeDoctorResponse[[]TimeDoctorUser], err error) {
	url, err := url.Parse(baseurl1_0 + "/users")
	utils.CheckError(err)

	query := url.Query()
	query.Add("company", companyId)
	query.Add("deleted", "false")
	query.Add("limit", "100")
	query.Add("token", token)

	url.RawQuery = query.Encode()

	req, err := http.Get(url.String())
	utils.CheckError(err)

	defer req.Body.Close()
	err = json.NewDecoder(req.Body).Decode(&response)
	utils.CheckError(err)

	return
}

func GenerateReportFromTimedoctor(token, company, user string, processImages bool, date time.Time) (data models.TimeDoctorReportForAnalysis, err error) {
	var wg sync.WaitGroup
	var imageData []TimeDoctorImageData
	var activityData []TimeDoctorActivity
	var tasks []string
  if processImages {
    wg.Add(1)
  }
	wg.Add(2)
	go func() {
		defer wg.Done()
		activityData, err = GetTimeuseData(token, company, user, date)
	}()
	go func() {
    if !processImages {
      return
    }
		defer wg.Done()
		imageData, err = GetActivityImageData(token, company, user, date)
	}()
	go func() {
		defer wg.Done()
		tasks, err = GetActiveTasks(token, company, user)
	}()

	wg.Wait()
	data = models.TimeDoctorReportForAnalysis{
		Activities: activityData,
		Images:     imageData,
		Tasks:      tasks,
		CreatedAt:  time.Now(),
	}

	return
}

func GetTimeuseData(token, company, user string, date time.Time) (data []TimeDoctorActivity, err error) {
	url, err := url.Parse(baseurl1_0 + "/activity/timeuse")
	if err != nil {
		return
	}

	start := "2023-11-01T18:30:00Z"
  // BeginningOfDay(date).UTC().Format(time.RFC3339)
	end := "2023-11-14T18:29:59Z"
  // EndOfDay(date).UTC().Format(time.RFC3339)

	query := url.Query()
	query.Add("company", company)
	query.Add("from", start)
	query.Add("to", end)
	query.Add("user", user)
	query.Add("token", token)
	query.Add("category-details", "true")
	query.Add("exclude-fields", "score")

	url.RawQuery = query.Encode()
	res, err := http.Get(url.String())

	if err != nil {
		return
	}
	defer res.Body.Close()

	activityData := &TimeDoctorResponse[[][]TimeDoctorActivity]{}
	err = json.NewDecoder(res.Body).Decode(&activityData)

  if len(activityData.Data) == 0 {
    return
  }
	data = activityData.Data[0]

	return
}

func GetActivityImageData(token, company, user string, date time.Time) (data []TimeDoctorImageData, err error) {
	url, err := url.Parse(baseurl1_0 + "/files")
	utils.CheckError(err)

	start := BeginningOfDay(date).UTC().Format(time.RFC3339)
	end := EndOfDay(date).UTC().Format(time.RFC3339)

	query := url.Query()
	query.Add("company", company)
	query.Add("user", user)
	query.Add("token", token)
	query.Add("filter[date]", start+"_"+end)

	url.RawQuery = query.Encode()

	res, err := http.Get(url.String())
	imageData := &TimeDoctorResponse[[]TimeDoctorImageData]{}
	err = json.NewDecoder(res.Body).Decode(&imageData)
	utils.CheckError(err)

	data = imageData.Data

	return
}

func GetActiveTasks(token, company, user string) (tasks []string, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	url, err := url.Parse(baseurl1_0 + "/stats/all/tasks")
	if err != nil {
		return
	}

	query := url.Query()
	query.Add("company", company)
	query.Add("user", user)
	query.Add("token", token)

	url.RawQuery = query.Encode()

	req, err := http.Get(url.String())
	res := &TimeDoctorResponse[[]map[string]any]{}
	err = json.NewDecoder(req.Body).Decode(&res)
	utils.CheckError(err)

	tasksToFetch := make(map[string]bool)

	for _, task := range res.Data {
		tasksToFetch[task["taskId"].(string)] = true
	}

	for id := range tasksToFetch {
		if id == "" {
			continue
		}
		url := baseurl1_0 + "/tasks/" + id + "?token=" + token + "&company=" + company
		req, err := http.Get(url)
		utils.CheckError(err)
		res := &TimeDoctorResponse[map[string]any]{}
		err = json.NewDecoder(req.Body).Decode(&res)
		tasks = append(tasks, res.Data["name"].(string))
	}

	return
}

func BeginningOfDay(date time.Time) time.Time {
  return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func EndOfDay(date time.Time) time.Time {
  return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
}
