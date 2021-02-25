//
// https://raw.githubusercontent.com/creekorful/mvnparser/master/parser.go
//
// MIT License
//
// Copyright (c) 2019 Aloïs Micard
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/xml"
	"io"
)

// Represent a POM file
type MavenProject struct {
	XMLName              xml.Name             `xml:"project"`
	ModelVersion         string               `xml:"modelVersion"`
	Parent               Parent               `xml:"parent"`
	GroupId              string               `xml:"groupId"`
	ArtifactId           string               `xml:"artifactId"`
	Version              string               `xml:"version"`
	Packaging            string               `xml:"packaging"`
	Name                 string               `xml:"name"`
	Repositories         []Repository         `xml:"repositories>repository"`
	Properties           Properties           `xml:"properties"`
	DependencyManagement DependencyManagement `xml:"dependencyManagement"`
	Dependencies         []Dependency         `xml:"dependencies>dependency"`
	Profiles             []Profile            `xml:"profiles"`
	Build                Build                `xml:"build"`
	PluginRepositories   []PluginRepository   `xml:"pluginRepositories>pluginRepository"`
	Modules              []string             `xml:"modules>module"`
}

// Represent the parent of the project
type Parent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// Represent a dependency of the project
type Dependency struct {
	XMLName    xml.Name    `xml:"dependency"`
	GroupId    string      `xml:"groupId"`
	ArtifactId string      `xml:"artifactId"`
	Version    string      `xml:"version"`
	Classifier string      `xml:"classifier"`
	Type       string      `xml:"type"`
	Scope      string      `xml:"scope"`
	Exclusions []Exclusion `xml:"exclusions>exclusion"`
}

// Represent an exclusion
type Exclusion struct {
	XMLName    xml.Name `xml:"exclusion"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

// Represent a repository
type Repository struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

type Profile struct {
	Id    string `xml:"id"`
	Build Build  `xml:"build"`
}

type Build struct {
	// todo: final name ?
	Plugins []Plugin `xml:"plugins>plugin"`
}

type Plugin struct {
	XMLName    xml.Name `xml:"plugin"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	//todo something like: Configuration map[string]string `xml:"configuration"`
	// todo executions
}

// Represent a pluginRepository
type PluginRepository struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
	Url  string `xml:"url"`
}

// Represent Properties
type Properties map[string]string

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*p = map[string]string{}
	for {
		key := ""
		value := ""
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			key = tokenType.Name.Local
			err := d.DecodeElement(&value, &start)
			if err != nil {
				return err
			}
			(*p)[key] = value
		}
	}
	return nil
}