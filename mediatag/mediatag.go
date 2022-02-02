package main

/*
This script looks up a video in the interactivepublisher database and returns an HTML video tag with the URL of the video in
*/

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guardian/new-encodings-endpoints/common"
	"html/template"
	"log"
	"strings"
)

var ops common.DynamoDbOps
var config common.Config
var mimeEquivelentsCache common.MimeEquivalentsCache

const HtmlTagTemplate = `<video preload='auto' id='video_{{.OctopusId}}' poster='{{.PosterURL}}'{{.ExtraArguments|attr}}>
  <source src='{{.Url}}' type='{{.Format}}'>
</video>`

type TemplateData struct {
	common.ContentResult
	ExtraArguments string
}

/*
templateHTML renders an html tag for the given found content that is injection-safe.

Arguments:

- foundContent - a non-NULL pointer to a ContentResult instance giving the content to build the tag for
- extraArguments - a string of extra arguments to put into the video tag

Returns:

- a string of the rendered html on success
- an error on failure.
*/
func templateHTML(foundContent *common.ContentResult, extraArguments string) (string, error) {
	//see https://stackoverflow.com/questions/14765395/why-am-i-seeing-zgotmplz-in-my-go-html-template-output
	extraFuncMap := template.FuncMap{
		//defines a "filter function" that marks the text as html-safe
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		//defines a "filter function" that marks the text as an HTML attribute
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
	}

	tmpl, err := template.New("html").Funcs(extraFuncMap).Parse(HtmlTagTemplate)
	if err != nil {
		log.Printf("ERROR Could not parse builtin html template \"%s\": %s", HtmlTagTemplate, err)
		return "", err
	}

	templateData := &TemplateData{
		ContentResult:  *foundContent,
		ExtraArguments: extraArguments,
	}

	wr := &strings.Builder{}
	err = tmpl.Execute(wr, templateData)
	if err != nil {
		log.Printf("ERROR Could not render found content %v to template \"%s\": %s", foundContent, HtmlTagTemplate, err)
		return "", err
	}
	return wr.String(), nil
}

func HandleEvent(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	foundContent, errResponse := common.FindContent(ctx, &event.QueryStringParameters, ops, config, mimeEquivelentsCache)
	if errResponse != nil {
		switch errResponse.StatusCode {
		case 404:
			return common.MakeResponseRaw(404, aws.String("No content found.\n"), "text/plain;charset=UTF-8"), nil
		default:
			return errResponse, nil
		}
	}

	extraArguments := ""
	if _, hasNoControls := (event.QueryStringParameters)["nocontrols"]; hasNoControls == false {
		extraArguments = extraArguments + " controls"
	}
	if _, hasAutoPlay := (event.QueryStringParameters)["autoplay"]; hasAutoPlay {
		extraArguments = extraArguments + " autoplay"
	}
	if _, hasLoop := (event.QueryStringParameters)["loop"]; hasLoop {
		extraArguments = extraArguments + " loop"
	}

	hTMLToReturn, err := templateHTML(foundContent, extraArguments)

	if err != nil {
		return common.MakeResponseJson(500, common.GenericErrorBody("Internal error, see logs")), nil
	}
	return common.MakeResponseRaw(200, &hTMLToReturn, "text/html;charset=UTF-8"), nil
}

func main() {
	var err error
	config, err = common.NewConfig()
	if err != nil {
		log.Printf("ERROR Could not initialise config: %s", err)
		panic("could not initialise config")
	}

	ops = common.NewDynamoDbOps(config)
	mimeEquivelentsCache, err = common.NewMimeEquivalentsCache(context.Background(), ops)
	if err != nil {
		log.Printf("ERROR Could not initialise mime equivalents: %s", err)
		panic("could not initialise MIME equivalents")
	}
	lambda.Start(HandleEvent)
}
