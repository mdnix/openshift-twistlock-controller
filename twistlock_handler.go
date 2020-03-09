package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
)

func getRolebinding(obj interface{}, action string) *Rolebinding {
	role := obj.(*rbacv1.RoleBinding)
	name := role.ObjectMeta.Name
	namespace := role.ObjectMeta.Namespace
	subjects := role.Subjects

	var rb *Rolebinding
	rb = &Rolebinding{
		Name:      name,
		Namespace: namespace,
		Action:    action,
	}

	for i := 0; i < len(subjects); i++ {
		group := subjects[i].Name
		kind := subjects[i].Kind

		isgroup, err := regexp.MatchString("^CN", group)
		if err != nil {
			panic(err.Error())
		}

		iskind, err := regexp.MatchString("^Group", kind)
		if err != nil {
			panic(err.Error())
		}

		if len(group) > 0 && iskind && isgroup {
			pattern := `CN=(?P<groupcn>.*?),`
			pathMetadata := regexp.MustCompile(pattern)
			matches := pathMetadata.FindStringSubmatch(group)
			var cn string

			for i, match := range matches {
				if i != 0 {
					cn = match
				}
			}
			rb.Group = append(rb.Group, group)
			rb.CN = append(rb.CN, cn)

			if strings.Contains(strings.ToLower(group), "admin") {
				rb.Role = "admin"
			} else {
				rb.Role = "devOps"
			}
		}
	}
	return rb
}

func parseCollection(twc TwistlockCollection) bytes.Buffer {
	var tmplBytes bytes.Buffer
	tmpl, err := template.ParseFiles(*configPath + "/twistlock-templates/collection.json")
	if err != nil {
		log.Panic(err)
	}
	tmpl.Execute(&tmplBytes, twc)
	return tmplBytes
}

func parseGroup(twg TwistlockGroup) bytes.Buffer {
	var tmplBytes bytes.Buffer
	tmpl, err := template.ParseFiles(*configPath + "/twistlock-templates/group.json")
	if err != nil {
		log.Panic(err)
	}
	tmpl.Execute(&tmplBytes, twg)
	return tmplBytes
}

func gettwAPI(endpoint string) []byte {
	var body []byte
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", twc.Host+endpoint, nil)
	if err != nil {
		logrus.Warnf("Unable to generate request to %s%s", twc.Host, endpoint)
	}

	req.SetBasicAuth(twc.User, twc.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "twistlock-controller")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("Unable to get json from %s%s. Error: %s", twc.Host, endpoint, err)
	}
	defer resp.Body.Close()
	logrus.Infof("Method: %s, Request: %s, Response: %s, Error: %s", req.Method, req.URL, resp.StatusCode, err)

	if resp.StatusCode == http.StatusOK {
		var err error
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Warn("Could not read bytes from body")
		}

	}
	logrus.Info("Data received")
	return body
}

func posttwAPI(endpoint string, payload string) int {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", twc.Host+endpoint, strings.NewReader(payload))
	if err != nil {
		logrus.Warnf("Unable to generate request to %s%s", twc.Host, endpoint)
	}
	req.SetBasicAuth(twc.User, twc.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "twistlock-controller")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("Could not post json to %s%s, Statuscode: %s, Error: %s", twc.Host, endpoint, resp.StatusCode, err)
	}
	logrus.Infof("Method: %s, Request: %s, Response: %s, Error: %s", req.Method, req.URL, resp.StatusCode, err)
	defer resp.Body.Close()
	logrus.Info("Data has been posted")
	return resp.StatusCode
}

func modifytwAPI(endpoint string, obj string, payload string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	safeobj := url.PathEscape(obj)

	req, err := http.NewRequest("PUT", twc.Host+endpoint+"/"+safeobj, strings.NewReader(payload))
	if err != nil {
		logrus.Warnf("Unable to generate request to %s%s%s%s", twc.Host, endpoint, "/", safeobj)
	}
	req.SetBasicAuth(twc.User, twc.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "twistlock-controller")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("Could not post json to %s%s, Statuscode: %s, Error: %s", twc.Host, endpoint, resp.StatusCode, err)
	}
	logrus.Infof("Method: %s, Request: %s, Response: %s, Error: %s", req.Method, req.URL, resp.StatusCode, err)
	defer resp.Body.Close()
	logrus.Info("Data has been modified")
}

func deletetwAPI(endpoint string, obj string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	safeobj := url.PathEscape(obj)
	req, err := http.NewRequest("DELETE", twc.Host+endpoint+"/"+safeobj, nil)
	if err != nil {
		logrus.Warnf("Unable to generate request to %s%s%s%s", twc.Host, endpoint, "/", safeobj)
	}
	req.SetBasicAuth(twc.User, twc.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "twistlock-controller")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("Could not delete obj %s, Statuscode: %s, Error: %s", obj, resp.StatusCode, err)
	}
	logrus.Infof("Method: %s, Request: %s, Response: %s, Error: %s", req.Method, req.URL, resp.StatusCode, err)
	defer resp.Body.Close()
	logrus.Info("Data has been deleted")
}

// Twistlock handler implements Handler interface,
type Twistlock struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (t *Twistlock) Init(c Config) error {

	return nil
}

// ObjectCreated sends events on object creation
func (t *Twistlock) ObjectCreated(obj interface{}) {
	etcdKey := fmt.Sprintf("%s/%s", obj.(*rbacv1.RoleBinding).Namespace, obj.(*rbacv1.RoleBinding).Name)
	etcdObj, err := json.Marshal(obj.(*rbacv1.RoleBinding))
	if err != nil {
		logrus.Warn(err)
	}

	_, err = kvPut(etcdKey, string(etcdObj))
	if err != nil {
		logrus.Warnf("Unable to put %s to etcd cluster", etcdKey)
	} else {
		logrus.Info("Rolebinding stored successfully on etcd cluster with key ", etcdKey)
	}

	role := getRolebinding(obj, "add")
	for _, cn := range role.CN {
		twcoll := TwistlockCollection{
			CN:        cn,
			Namespace: role.Namespace,
		}
		if role.Role == "devOps" {
			var collExists *bool
			collBytes := gettwAPI(twcollAPI)
			if len(collBytes) > 0 {
				var collections []CollectionAPI
				err := json.Unmarshal(collBytes, &collections)
				if err != nil {
					logrus.Warn(err)
				}
				for i := 0; i < len(collections); i++ {
					if collections[i].Name == twcoll.CN {
						logrus.Infof("Collection %s already exists", collections[i])
						if !sliceContains(collections[i].Namespaces, twcoll.Namespace) {
							t := true
							collExists = &t
							logrus.Infof("Adding namespace %s to collection %s", twcoll.Namespace, collections[i].Name)
							collections[i].Namespaces = append(collections[i].Namespaces, twcoll.Namespace)
							data, _ := json.Marshal(collections[i])
							modifytwAPI(twcollAPI, collections[i].Name, string(data))
							break
						} else {
							t := true
							collExists = &t
							break
						}
					} else {
						f := false
						collExists = &f
					}
				}
			}
			if !*collExists {
				logrus.Infof("Creating Collection %s", twcoll.CN)
				twObj := parseCollection(twcoll)
				jsonObj := twObj.String()
				s := posttwAPI(twcollAPI, jsonObj)
				if s >= 200 && s <= 299 {
					logrus.Info("Collection posted successfully")
				} else {
					logrus.Info("Unable to post collection")
				}
			}

		}
	}
	for i, group := range role.Group {
		twgroup := TwistlockGroup{
			CN:    role.CN[i],
			Group: group,
			Role:  role.Role,
		}
		if twgroup.Role == "devOps" {
			var grpExits *bool
			groupBytes := gettwAPI(twgrpAPI)
			if len(groupBytes) > 0 {
				var groups []GroupAPI
				err := json.Unmarshal(groupBytes, &groups)
				if err != nil {
					logrus.Warn(err)
				}
				for i := 0; i < len(groups); i++ {
					if groups[i].GroupName == twgroup.CN {
						t := true
						grpExits = &t
						logrus.Infof("Group %s already exists", groups[i])
						break
					} else {
						f := false
						grpExits = &f
					}
				}
			}
			if !*grpExits {
				logrus.Infof("Creating Group %s", twgroup.CN)
				twObj := parseGroup(twgroup)
				jsonObj := twObj.String()
				s := posttwAPI(twgrpAPI, jsonObj)
				if s >= 200 && s <= 299 {
					logrus.Info("Group posted successfully")
				} else {
					logrus.Info("Unable to post group")
				}
			}

		}
	}
}

// ObjectUpdated sends events on object updation
func (t *Twistlock) ObjectUpdated(obj interface{}) {
	var add []string
	var del []string
	newRole := getRolebinding(obj.(Event).newObj, "update")
	oldRole := getRolebinding(obj.(Event).oldObj, "update")

	etcdKey := fmt.Sprintf("%s/%s", newRole.Namespace, newRole.Name)
	etcdObj, err := json.Marshal(newRole)
	if err != nil {
		logrus.Warn(err)
	}

	_, err = kvPut(etcdKey, string(etcdObj))
	if err != nil {
		logrus.Warnf("Unable to put %s to etcd cluster", etcdKey)
	} else {
		logrus.Info("Rolebinding updated successfully on etcd cluster with key ", etcdKey)
	}

	for _, cn := range oldRole.CN {
		c := sliceContains(newRole.CN, cn)
		if !c {
			logrus.Info("This group got deleted: ", cn)
			del = append(del, cn)
		}
	}
	for _, cn := range newRole.CN {
		c := sliceContains(oldRole.CN, cn)
		if !c {
			logrus.Info("This group got added: ", cn)
			add = append(add, cn)
		}
	}

	if len(del) > 0 {
		for _, cn := range del {
			if newRole.Role == "devOps" {
				collBytes := gettwAPI(twcollAPI)
				if len(collBytes) > 0 {
					var collections []CollectionAPI
					err := json.Unmarshal(collBytes, &collections)
					if err != nil {
						logrus.Warn(err)
					}
					for i := 0; i < len(collections); i++ {
						if collections[i].Name == cn {
							ns := len(collections[i].Namespaces)

							if ns <= 1 && sliceContains(collections[i].Namespaces, newRole.Namespace) {
								logrus.Infof("%s is the only namespace in collection %s", newRole.Namespace, collections[i].Name)
								logrus.Infof("Deleting Group %s", cn)
								deletetwAPI(twgrpAPI, cn)
								logrus.Infof("Deleting collection %s", cn)
								deletetwAPI(twcollAPI, cn)
							} else if ns > 1 && sliceContains(collections[i].Namespaces, newRole.Namespace) {
								logrus.Infof("Removing namespace %s from collection %s", newRole.Namespace, collections[i].Name)
								collections[i].Namespaces = sliceRemove(collections[i].Namespaces, newRole.Namespace)
								data, _ := json.Marshal(collections[i])
								modifytwAPI(twcollAPI, collections[i].Name, string(data))
							}
						}
					}

				}
			}
		}
	}

	if len(add) > 0 {
		for _, cn := range add {
			twcoll := TwistlockCollection{
				CN:        cn,
				Namespace: newRole.Namespace,
			}
			if newRole.Role == "devOps" {
				var collExists *bool
				collBytes := gettwAPI(twcollAPI)
				if len(collBytes) > 0 {
					var collections []CollectionAPI
					err := json.Unmarshal(collBytes, &collections)
					if err != nil {
						logrus.Warn(err)
					}
					for i := 0; i < len(collections); i++ {
						if collections[i].Name == twcoll.CN {
							logrus.Infof("Collection %s already exists", collections[i])
							if !sliceContains(collections[i].Namespaces, twcoll.Namespace) {
								t := true
								collExists = &t
								logrus.Infof("Adding namespace %s to collection %s", twcoll.Namespace, collections[i].Name)
								collections[i].Namespaces = append(collections[i].Namespaces, twcoll.Namespace)
								data, _ := json.Marshal(collections[i])
								modifytwAPI(twcollAPI, collections[i].Name, string(data))
								break
							} else {
								t := true
								collExists = &t
								break
							}
						} else {
							f := false
							collExists = &f
						}
					}
				}
				if !*collExists {
					logrus.Infof("Creating Collection %s", twcoll.CN)
					twObj := parseCollection(twcoll)
					jsonObj := twObj.String()
					s := posttwAPI(twcollAPI, jsonObj)
					if s >= 200 && s <= 299 {
						logrus.Info("Collection posted successfully")
					} else {
						logrus.Info("Unable to post collection")
					}
				}

			}
		}
		for _, group := range add {
			twgroup := TwistlockGroup{
				CN:   group,
				Role: newRole.Role,
			}
			if twgroup.Role == "devOps" {
				var grpExits *bool
				groupBytes := gettwAPI(twgrpAPI)
				if len(groupBytes) > 0 {
					var groups []GroupAPI
					err := json.Unmarshal(groupBytes, &groups)
					if err != nil {
						logrus.Warn(err)
					}
					for i := 0; i < len(groups); i++ {
						if groups[i].GroupName == twgroup.CN {
							t := true
							grpExits = &t
							logrus.Infof("Group %s already exists", groups[i])
							break
						} else {
							f := false
							grpExits = &f
						}
					}
				}
				if !*grpExits {
					logrus.Infof("Creating Group %s", twgroup.CN)
					twObj := parseGroup(twgroup)
					jsonObj := twObj.String()
					s := posttwAPI(twgrpAPI, jsonObj)
					if s >= 200 && s <= 299 {
						logrus.Info("Group posted successfully")
					} else {
						logrus.Info("Unable to post group")
					}
				}

			}
		}
	}
}

// ObjectDeleted sends events on object deletion
func (t *Twistlock) ObjectDeleted(obj interface{}) {
	etcdKey := obj.(string)
	etcdObj, err := kvGet(etcdKey)
	if err != nil {
		logrus.Warnf("Unable to get %s from etcd cluster", etcdKey)
	}

	var rb *rbacv1.RoleBinding
	err = json.Unmarshal(etcdObj.Kvs[0].Value, &rb)
	if err != nil {
		logrus.Warn("Unable to unmarshal rolebinding: ", err)
	}
	role := getRolebinding(rb, "delete")

	for _, cn := range role.CN {
		twcoll := TwistlockCollection{
			CN:        cn,
			Namespace: role.Namespace,
		}
		if role.Role == "devOps" {
			collBytes := gettwAPI(twcollAPI)
			if len(collBytes) > 0 {
				var collections []CollectionAPI
				err := json.Unmarshal(collBytes, &collections)
				if err != nil {
					logrus.Warn(err)
				}
				for i := 0; i < len(collections); i++ {
					if collections[i].Name == twcoll.CN {
						ns := len(collections[i].Namespaces)

						if ns <= 1 && sliceContains(collections[i].Namespaces, twcoll.Namespace) {
							logrus.Infof("%s is the only namespace in collection %s", twcoll.Namespace, collections[i].Name)
							logrus.Infof("Deleting Group %s", twcoll.CN)
							deletetwAPI(twgrpAPI, twcoll.CN)
							logrus.Infof("Deleting collection %s", twcoll.CN)
							deletetwAPI(twcollAPI, twcoll.CN)
						} else if ns > 1 && sliceContains(collections[i].Namespaces, twcoll.Namespace) {
							logrus.Infof("Removing namespace %s from collection %s", twcoll.Namespace, collections[i].Name)
							collections[i].Namespaces = sliceRemove(collections[i].Namespaces, twcoll.Namespace)
							data, _ := json.Marshal(collections[i])
							modifytwAPI(twcollAPI, collections[i].Name, string(data))
						}
						_, err := kvDel(etcdKey)
						if err != nil {
							logrus.Warnf("Unable to delete %s from etcd cluster", etcdKey)
						} else {
							logrus.Info("Rolebinding deleted successfully from etcd cluster with key ", etcdKey)
						}
					}
				}

			}
		}
	}
}
