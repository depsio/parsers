package main

import (
	"fmt"
	"os"

	"encoding/json"

	glidecfg "github.com/Masterminds/glide/cfg"
)

type FileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Request struct {
	Files []FileRequest `json:"files"`
}

type Specifier struct {
	Operator string `json:"operator,omitempty"`
	Version  string `json:"version"`
}

type Dependency struct {
	Name       string                 `json:"name"`
	Specifiers []Specifier            `json:"specifiers"`
	Group      string                 `json:"group,omitempty"`
	Extras     map[string]interface{} `json:"extras"`
}

type FileResponse struct {
	Name         string       `json:"name"`
	Dependencies []Dependency `json:"dependencies"`
	Error        string       `json:"error,omitempty"`
}

type Response struct {
	Files []FileResponse `json:"files"`
}

func main() {
	var req Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "can't parse input JSON: %s\n", err)
		os.Exit(1)
	}

	var resp Response
	for _, fileReq := range req.Files {
		fileResp := FileResponse{
			Name: fileReq.Name,
		}

		cfg, err := glidecfg.ConfigFromYaml([]byte(fileReq.Content))
		if err == nil {
			for _, src := range []struct {
				group string
				deps  glidecfg.Dependencies
			}{
				{"", cfg.Imports},
				{"development", cfg.DevImports},
			} {
				for _, dep := range src.deps {
					depResp := Dependency{
						Name:   dep.Name,
						Extras: make(map[string]interface{}),
						Group:  src.group,
					}

					if dep.Reference != "" {
						depResp.Specifiers = []Specifier{
							{"", dep.Reference},
						}
						depResp.Extras["reference"] = dep.Reference
					}

					if dep.Repository != "" {
						depResp.Extras["repo"] = dep.Repository
					}

					if dep.VcsType != "" {
						depResp.Extras["vcs"] = dep.VcsType
					}

					if len(dep.Subpackages) > 0 {
						depResp.Extras["subpackages"] = dep.Subpackages
					}

					if len(dep.Os) > 0 {
						depResp.Extras["os"] = dep.Os
					}

					if len(dep.Arch) > 0 {
						depResp.Extras["arch"] = dep.Arch
					}

					fileResp.Dependencies = append(fileResp.Dependencies, depResp)
				}
			}
		} else {
			fileResp.Error = fmt.Sprintf("can't parse input YAML: %s", err)
		}

		resp.Files = append(resp.Files, fileResp)
	}

	json.NewEncoder(os.Stdout).Encode(resp)
}
