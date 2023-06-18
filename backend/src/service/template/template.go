// indece Monitor
// Copyright (C) 2023 indece UG (haftungsbeschränkt)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License or any
// later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package template

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	texttemplate "text/template"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/assets"
	"github.com/indece-official/monitor/backend/src/model"
)

const ServiceName = "template"

type ExecutionData struct {
	Data map[string]interface{}
}

type ITemplate interface {
	Execute(io.Writer, interface{}) error
}

type TemplateType string

const (
	TemplateTypeStatusChanged TemplateType = "status_changed"
)

var TemplateTypes = []TemplateType{
	TemplateTypeStatusChanged,
}

type TemplatePart string

const (
	TemplatePartEmailBodyHTML TemplatePart = "email_body.html"
	TemplatePartEmailBodyText TemplatePart = "email_body.txt"
	TemplatePartEmailSubject  TemplatePart = "email_subject.txt"
)

var TemplatePartsEmail = []TemplatePart{
	TemplatePartEmailBodyHTML,
	TemplatePartEmailBodyText,
	TemplatePartEmailSubject,
}

type IService interface {
	gousu.IService

	Generate(locale model.Locale, templateType TemplateType, templatePart TemplatePart, params map[string]interface{}) (string, error)
}

type Service struct {
	log       *logger.Log
	templates map[model.Locale]map[TemplateType]map[TemplatePart]ITemplate
}

var _ (IService) = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) getFilename(locale model.Locale, templateType TemplateType, templatePart TemplatePart) string {
	return fmt.Sprintf("templates/%s/%s/%s", locale, templateType, templatePart)
}

func (s *Service) loadTemplate(locale model.Locale, templateType TemplateType, templatePart TemplatePart) error {
	filename := s.getFilename(locale, templateType, templatePart)

	templateFile, err := assets.Assets.Open(filename)
	if err != nil {
		return fmt.Errorf("can't open template file: %s", err)
	}
	defer templateFile.Close()

	templateStr, err := io.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("can't read template file: %s", err)
	}

	var tpl ITemplate

	switch templatePart {
	case TemplatePartEmailBodyHTML:
		tpl, err = htmltemplate.New(filename).Parse(string(templateStr))
		if err != nil {
			return fmt.Errorf("can't load html template: %s", err)
		}
	default:
		tpl, err = texttemplate.New(filename).Parse(string(templateStr))
		if err != nil {
			return fmt.Errorf("can't load text template: %s", err)
		}
	}

	if _, ok := s.templates[locale]; !ok {
		s.templates[locale] = map[TemplateType]map[TemplatePart]ITemplate{}
	}

	if _, ok := s.templates[locale][templateType]; !ok {
		s.templates[locale][templateType] = map[TemplatePart]ITemplate{}
	}

	s.templates[locale][templateType][templatePart] = tpl

	return nil
}

func (s *Service) Start() error {
	s.templates = map[model.Locale]map[TemplateType]map[TemplatePart]ITemplate{}

	for _, locale := range model.Locales {
		for _, templateType := range TemplateTypes {
			for _, templatePart := range TemplatePartsEmail {
				err := s.loadTemplate(locale, templateType, templatePart)
				if err != nil {
					return fmt.Errorf("error loading template %s:%s:%s: %s", locale, templateType, templatePart, err)
				}
			}
		}
	}

	return nil
}

func (s *Service) Health() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Generate(
	locale model.Locale,
	templateType TemplateType,
	templatePart TemplatePart,
	params map[string]interface{},
) (string, error) {
	var buf bytes.Buffer

	if _, ok := s.templates[locale]; !ok {
		return "", fmt.Errorf("template not found: unknown locale '%s'", locale)
	}

	if _, ok := s.templates[locale][templateType]; !ok {
		return "", fmt.Errorf("template not found: unknown template type '%s'", templateType)
	}

	if _, ok := s.templates[locale][templateType][templatePart]; !ok {
		return "", fmt.Errorf("template not found: unknown template part '%s'", templatePart)
	}

	data := &ExecutionData{
		Data: params,
	}

	err := s.templates[locale][templateType][templatePart].Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("can't execute template: %s", err)
	}

	return buf.String(), nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
