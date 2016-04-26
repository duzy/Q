package main

import (
        "github.com/duzy/W"
        "encoding/json"
        "net/http"
        "strconv"
        "strings"
        "fmt"
        "log"
        //"io"
)

var (
        quizAnswerNameFmt = "Q.A.%v"
        quizNameFmt = "%v.quiz.tmpl"
        //quizList = []string{ "1", "2", "3", "4", "5", "6" }
)

func init() {
        a := []string{
                "temp/*.tmpl",
                "temp/Q/*.quiz.tmpl",
                "temp/header",
                "temp/footer",
                "temp/page",
                "temp/view",
        }
        W.
                Delims("{{{", "}}}").
                MustGlob(a...)
}

var (
        users = map[string]*user{
                "a": &user{
                        answers: make(map[string]*answer),
                        firstName: "First",
                        lastName: "Last",
                        email: "a",
                        pass: "a",
                        online: false,
                        quizNum: 0,
                },
        }
)

type user struct {
        answers map[string]*answer
        email, pass string
        firstName, lastName string
        online bool
        quizNum int
}

type answer struct {
        a, x []string
}

func (a *answer) A() string { return strings.Join(a.a, ", ") }
func (a *answer) X() string { return strings.Join(a.x, ", ") }
func (a *answer) Check() bool { return a.isCorrect() }
func (a *answer) isCorrect() bool {
        for _, s := range a.a {
                found := false
                for _, x := range a.x {
                        if strings.EqualFold(s, x) {
                                found = true; break
                        }
                }
                if !found {
                        return false
                }
        }
        return 0 < len(a.a) && len(a.a) == len(a.x)
}

func dealTest(dc *W.DealContext) {
        dc.Set("Data", "this is a test")
}

func dealHome(dc *W.DealContext) {
}

func dealLogin(dc *W.DealContext) {
}

func dealRegister(dc *W.DealContext) {
}

// Note: It's better using RESTful and JSON instead.
func dealExam(dc *W.DealContext) {
        if err := dc.R.ParseForm(); err != nil {
                log.Fatalf("ParseForm: %v", err)
        }
        
        // TODO: check user TOKEN

        //log.Printf("Exam: %v, %v", dc.R.Form, dc.R.Cookies())

        var (
                token, _ = dc.R.Cookie("token")
                quizNum = 1
                quizName, name string
        )

        if token == nil || len(token.Value) == 0 || !strings.HasPrefix(token.Value, "user-") {
                dc.Set("Error", "unauthorized request")
                dc.Name = "error.tmpl"
                return
        } else {
                name = token.Value[5:]
        }
        
        if u, ok := users[name]; !ok || u == nil {
                dc.Set("Error", "invalid token")
                dc.Name = "error.tmpl"
        } else {
                if s := dc.R.Form.Get("n"); s != "" {
                        if i, err := strconv.Atoi(s); err == nil && 0 < i {
                                quizNum = i
                        }
                }

                quizName = fmt.Sprintf(quizNameFmt, quizNum);
                answered, exists := len(u.answers), W.CanExpand(quizName)
                
                //log.Printf("Exam: %v (exists: %v (%v)) (answered: %v)", quizName, exists, dc.Name, answered)
                
                if exists {
                        dc.Set("N", quizNum)
                        dc.Set("Q", quizName)
                        dc.Name = "exam.tmpl"
                } else if 0 < quizNum && answered+1 == quizNum {
                        correct := 0
                        for _, a := range u.answers {
                                if a.isCorrect() {
                                        correct++
                                }
                        }
                        dc.Set("Score", int(float32(correct)/float32(answered) * 100))
                        dc.Set("Answers", u.answers)
                        dc.Name = "results.tmpl"
                } else {
                        dc.Set("N", quizNum)
                        dc.Name = "exam_missing.tmpl"
                }
        }
}

func dealResults(dc *W.DealContext) {
}

func dealMarks(dc *W.DealContext) {
}

func dealLogout(dc *W.DealContext) {
}

func writeJSONContent(w http.ResponseWriter, m map[string]interface{}) {
        // application/json, application/javascript
        w.Header().Set("Content-Type", "application/json")
        if b, err := json.Marshal(m); err == nil {
                w.Write(b)
        } else {
                log.Fatalf("Marshal: %v", err)
        }
}

func handleUserRequest(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
                log.Fatalf("ParseForm: %v", err)
        }

        var res = make(map[string]interface{})
        
        log.Printf("User: %v", r.Form)

        switch r.Form.Get("method") {
        case "register": userRegister(w, r, res)
        case "login":    userLogin(w, r, res)
        case "logout":   userLogout(w, r, res)
        default:
                res["message"] = "unsupported method"
                res["code"] = -1
        }

        writeJSONContent(w, res)
}

func userRegister(w http.ResponseWriter, r *http.Request, res map[string]interface{}) {
        var (
                firstName = r.Form.Get("firstName")
                lastName = r.Form.Get("lastName")
                email = r.Form.Get("email")
                pass = r.Form.Get("pass")
        )
        if u, ok := users[email]; ok && u != nil {
                res["message"] = "user already exists"
                res["code"] = -2
        } else {
                users[email] = &user{
                        answers: make(map[string]*answer),
                        firstName: firstName,
                        lastName: lastName,
                        email: email,
                        pass: pass,
                        online: true,
                        quizNum: 0,
                }
                res["message"] = "okay"
                res["code"] = 0
                res["quiz"] = 0
                res["token"] = fmt.Sprintf("user-%s", email) // TODO: auth token
        }
}

func userLogin(w http.ResponseWriter, r *http.Request, res map[string]interface{}) {
        var (
                email = r.Form.Get("email")
                pass = r.Form.Get("pass")
        )
        if u, ok := users[email]; ok && u != nil {
                if pass == u.pass {
                        u.online = true
                        res["code"] = 0
                        res["message"] = "okay"
                        res["quiz"] = u.quizNum
                        res["token"] = fmt.Sprintf("user-%s", email) // TODO: auth token
                } else {
                        res["code"] = -3
                        res["message"] = "invalid password"
                }
        } else {
                res["code"] = -3
                res["message"] = "no such user"
        }
}

func userLogout(w http.ResponseWriter, r *http.Request, res map[string]interface{}) {
        var (
                email = r.Form.Get("email")
                //pass = r.Form.Get("pass")
        )
        if u, ok := users[email]; ok && u != nil {
                u.online = false
                res["code"] = 0
                res["message"] = "okay"
        } else {
                res["code"] = -5
                res["message"] = "no such user"
        }
}

func handleTakeRequest(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
                log.Fatalf("ParseForm: %v", err)
        }

        //log.Printf("Take: %v", r.Form)
        
        var (
                res = make(map[string]interface{})
                ans = r.Form.Get("a")
                token = r.Form.Get("t")
                num = r.Form.Get("n")
        )

        defer writeJSONContent(w, res)
        
        if len(token) == 0 || !strings.HasPrefix(token, "user-") {
                res["message"] = "unauthorized request"
                res["code"] = -6
                return
        }
        
        name := token[5:]
        
        if u, ok := users[name]; !ok || u == nil {
                res["message"] = "no such user"
                res["code"] = -7
        } else {
                s, a := fmt.Sprintf("%v", num), new(answer)
                o := &struct{ A *[]string }{ &a.a }
                if err := json.Unmarshal([]byte(ans), o); err != nil {
                        log.Printf("Unmarshal: %v (%s)", err, ans)
                        res["message"] = "invalid answer"
                        res["code"] = -8
                } else {
                        s = fmt.Sprintf(quizAnswerNameFmt, s)
                        a.x = strings.Split(W.MustLoadString(s), ";")
                        for i, s := range a.x {
                                a.x[i] = strings.TrimSpace(s)
                        }
                        log.Printf("Answer: %v (correct=%v) (%s) (%s)", a, a.isCorrect(), s, ans)
                        u.answers[s] = a
                        res["message"] = "ok"
                        res["code"] = 0
                }
        }
}

func main() {
        assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
        http.Handle("/assets/",         assets)
        http.Handle("/",                W.Deal("home.tmpl",     dealHome))
        http.Handle("/test",            W.Deal("test.tmpl",     dealTest))
        http.Handle("/ok",              W.Deal("ok.tmpl",       dealTest))
        http.Handle("/v/register",      W.Deal("reg.tmpl",      dealRegister))
        http.Handle("/v/login",         W.Deal("login.tmpl",    dealLogin))
        http.Handle("/v/logout",        W.Deal("logout.tmpl",   dealLogout))
        http.Handle("/v/exam",          W.Deal("exam.tmpl",     dealExam))
        http.Handle("/v/results",       W.Deal("results.tmpl",  dealResults))
        http.Handle("/v/marks",         W.Deal("marks.tmpl",    dealMarks))
        http.HandleFunc("/user",        handleUserRequest)
        http.HandleFunc("/take",        handleTakeRequest)
        if err := http.ListenAndServe(":8088", nil); err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
