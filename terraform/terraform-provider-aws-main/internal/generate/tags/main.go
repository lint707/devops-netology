//go:build generate
// +build generate

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	v1 "github.com/hashicorp/terraform-provider-aws/internal/generate/tags/templates/v1"
	v2 "github.com/hashicorp/terraform-provider-aws/internal/generate/tags/templates/v2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	filename = `tags_gen.go`

	sdkV1 = 1
	sdkV2 = 2
)

var (
	getTag             = flag.Bool("GetTag", false, "whether to generate GetTag")
	listTags           = flag.Bool("ListTags", false, "whether to generate ListTags")
	serviceTagsMap     = flag.Bool("ServiceTagsMap", false, "whether to generate service tags for map")
	serviceTagsSlice   = flag.Bool("ServiceTagsSlice", false, "whether to generate service tags for slice")
	untagInNeedTagType = flag.Bool("UntagInNeedTagType", false, "whether Untag input needs tag type")
	updateTags         = flag.Bool("UpdateTags", false, "whether to generate UpdateTags")

	listTagsInFiltIDName  = flag.String("ListTagsInFiltIDName", "", "listTagsInFiltIDName")
	listTagsInIDElem      = flag.String("ListTagsInIDElem", "ResourceArn", "listTagsInIDElem")
	listTagsInIDNeedSlice = flag.String("ListTagsInIDNeedSlice", "", "listTagsInIDNeedSlice")
	listTagsOp            = flag.String("ListTagsOp", "ListTagsForResource", "listTagsOp")
	listTagsOutTagsElem   = flag.String("ListTagsOutTagsElem", "Tags", "listTagsOutTagsElem")
	tagInCustomVal        = flag.String("TagInCustomVal", "", "tagInCustomVal")
	tagInIDElem           = flag.String("TagInIDElem", "ResourceArn", "tagInIDElem")
	tagInIDNeedSlice      = flag.String("TagInIDNeedSlice", "", "tagInIDNeedSlice")
	tagInTagsElem         = flag.String("TagInTagsElem", "Tags", "tagInTagsElem")
	tagKeyType            = flag.String("TagKeyType", "", "tagKeyType")
	tagOp                 = flag.String("TagOp", "TagResource", "tagOp")
	tagOpBatchSize        = flag.String("TagOpBatchSize", "", "tagOpBatchSize")
	tagResTypeElem        = flag.String("TagResTypeElem", "", "tagResTypeElem")
	tagType               = flag.String("TagType", "Tag", "tagType")
	tagType2              = flag.String("TagType2", "", "tagType")
	TagTypeAddBoolElem    = flag.String("TagTypeAddBoolElem", "", "TagTypeAddBoolElem")
	tagTypeIDElem         = flag.String("TagTypeIDElem", "", "tagTypeIDElem")
	tagTypeKeyElem        = flag.String("TagTypeKeyElem", "Key", "tagTypeKeyElem")
	tagTypeValElem        = flag.String("TagTypeValElem", "Value", "tagTypeValElem")
	untagInCustomVal      = flag.String("UntagInCustomVal", "", "untagInCustomVal")
	untagInNeedTagKeyType = flag.String("UntagInNeedTagKeyType", "", "untagInNeedTagKeyType")
	untagInTagsElem       = flag.String("UntagInTagsElem", "TagKeys", "untagInTagsElem")
	untagOp               = flag.String("UntagOp", "UntagResource", "untagOp")

	parentNotFoundErrCode = flag.String("ParentNotFoundErrCode", "", "Parent 'NotFound' Error Code")
	parentNotFoundErrMsg  = flag.String("ParentNotFoundErrMsg", "", "Parent 'NotFound' Error Message")

	sdkVersion   = flag.Int("AWSSDKVersion", sdkV1, "Version of the AWS SDK Go to use i.e. 1 or 2")
	kvtValues    = flag.Bool("KVTValues", false, "Whether KVT string map is of string pointers")
	skipTypesImp = flag.Bool("SkipTypesImp", false, "Whether to skip importing types")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tmain.go [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

type TemplateBody struct {
	getTag           string
	header           string
	listTags         string
	serviceTagsMap   string
	serviceTagsSlice string
	updateTags       string
}

func NewTemplateBody(version int, kvtValues bool) *TemplateBody {
	switch version {
	case sdkV1:
		return &TemplateBody{
			"\n" + v1.GetTagBody,
			v1.HeaderBody,
			"\n" + v1.ListTagsBody,
			"\n" + v1.ServiceTagsMapBody,
			"\n" + v1.ServiceTagsSliceBody,
			"\n" + v1.UpdateTagsBody,
		}
	case sdkV2:
		if kvtValues {
			return &TemplateBody{
				"\n" + v2.GetTagBody,
				v2.HeaderBody,
				"\n" + v2.ListTagsBody,
				"\n" + v2.ServiceTagsValueMapBody,
				"\n" + v2.ServiceTagsSliceBody,
				"\n" + v2.UpdateTagsBody,
			}
		}
		return &TemplateBody{
			"\n" + v2.GetTagBody,
			v2.HeaderBody,
			"\n" + v2.ListTagsBody,
			"\n" + v2.ServiceTagsMapBody,
			"\n" + v2.ServiceTagsSliceBody,
			"\n" + v2.UpdateTagsBody,
		}
	default:
		return nil
	}
}

type TemplateData struct {
	AWSService             string
	AWSServiceIfacePackage string
	ClientType             string
	ServicePackage         string

	ListTagsInFiltIDName    string
	ListTagsInIDElem        string
	ListTagsInIDNeedSlice   string
	ListTagsOp              string
	ListTagsOutTagsElem     string
	ParentNotFoundErrCode   string
	ParentNotFoundErrMsg    string
	RetryCreateOnNotFound   string
	TagInCustomVal          string
	TagInIDElem             string
	TagInIDNeedSlice        string
	TagInTagsElem           string
	TagKeyType              string
	TagOp                   string
	TagOpBatchSize          string
	TagPackage              string
	TagResTypeElem          string
	TagType                 string
	TagType2                string
	TagTypeAddBoolElem      string
	TagTypeAddBoolElemSnake string
	TagTypeIDElem           string
	TagTypeKeyElem          string
	TagTypeValElem          string
	UntagInCustomVal        string
	UntagInNeedTagKeyType   string
	UntagInNeedTagType      bool
	UntagInTagsElem         string
	UntagOp                 string

	// The following are specific to writing import paths in the `headerBody`;
	// to include the package, set the corresponding field's value to true
	ContextPkg      bool
	FmtPkg          bool
	HelperSchemaPkg bool
	SkipTypesImp    bool
	StrConvPkg      bool
	TfResourcePkg   bool
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *sdkVersion != sdkV1 && *sdkVersion != sdkV2 {
		log.Fatalf("AWS SDK Go Version %d not supported", *sdkVersion)
	}

	servicePackage := os.Getenv("GOPACKAGE")
	awsPkg, err := names.AWSGoPackage(servicePackage, *sdkVersion)

	if err != nil {
		log.Fatalf("encountered: %s", err)
	}

	var awsIntfPkg string
	if *sdkVersion == sdkV1 && (*getTag || *listTags || *updateTags) {
		awsIntfPkg = fmt.Sprintf("%[1]s/%[1]siface", awsPkg)
	}

	clientTypeName, err := names.AWSGoClientTypeName(servicePackage, *sdkVersion)

	if err != nil {
		log.Fatalf("encountered: %s", err)
	}

	var clientType string
	if *sdkVersion == sdkV1 {
		clientType = fmt.Sprintf("%siface.%sAPI", awsPkg, clientTypeName)
	} else {
		clientType = fmt.Sprintf("*%s.%s", awsPkg, clientTypeName)
	}

	tagPackage := awsPkg

	if tagPackage == "wafregional" {
		tagPackage = "waf"
		if *sdkVersion == sdkV1 {
			awsPkg = ""
		}
	}

	templateData := TemplateData{
		AWSService:             awsPkg,
		AWSServiceIfacePackage: awsIntfPkg,
		ClientType:             clientType,
		ServicePackage:         servicePackage,

		ContextPkg:      *sdkVersion == sdkV2 || (*getTag || *listTags || *updateTags),
		FmtPkg:          *updateTags,
		HelperSchemaPkg: awsPkg == "autoscaling",
		SkipTypesImp:    *skipTypesImp,
		StrConvPkg:      awsPkg == "autoscaling",
		TfResourcePkg:   *getTag,

		ListTagsInFiltIDName:    *listTagsInFiltIDName,
		ListTagsInIDElem:        *listTagsInIDElem,
		ListTagsInIDNeedSlice:   *listTagsInIDNeedSlice,
		ListTagsOp:              *listTagsOp,
		ListTagsOutTagsElem:     *listTagsOutTagsElem,
		ParentNotFoundErrCode:   *parentNotFoundErrCode,
		ParentNotFoundErrMsg:    *parentNotFoundErrMsg,
		TagInCustomVal:          *tagInCustomVal,
		TagInIDElem:             *tagInIDElem,
		TagInIDNeedSlice:        *tagInIDNeedSlice,
		TagInTagsElem:           *tagInTagsElem,
		TagKeyType:              *tagKeyType,
		TagOp:                   *tagOp,
		TagOpBatchSize:          *tagOpBatchSize,
		TagPackage:              tagPackage,
		TagResTypeElem:          *tagResTypeElem,
		TagType:                 *tagType,
		TagType2:                *tagType2,
		TagTypeAddBoolElem:      *TagTypeAddBoolElem,
		TagTypeAddBoolElemSnake: ToSnakeCase(*TagTypeAddBoolElem),
		TagTypeIDElem:           *tagTypeIDElem,
		TagTypeKeyElem:          *tagTypeKeyElem,
		TagTypeValElem:          *tagTypeValElem,
		UntagInCustomVal:        *untagInCustomVal,
		UntagInNeedTagKeyType:   *untagInNeedTagKeyType,
		UntagInNeedTagType:      *untagInNeedTagType,
		UntagInTagsElem:         *untagInTagsElem,
		UntagOp:                 *untagOp,
	}

	templateBody := NewTemplateBody(*sdkVersion, *kvtValues)

	if *getTag || *listTags || *serviceTagsMap || *serviceTagsSlice || *updateTags {
		// If you intend to only generate Tags and KeyValueTags helper methods,
		// the corresponding aws-sdk-go	service package does not need to be imported
		if !*getTag && !*listTags && !*serviceTagsSlice && !*updateTags {
			templateData.AWSService = ""
			templateData.TagPackage = ""
		}
		writeTemplate(templateBody.header, "header", templateData)
	}

	if *getTag {
		writeTemplate(templateBody.getTag, "gettag", templateData)
	}

	if *listTags {
		writeTemplate(templateBody.listTags, "listtags", templateData)
	}

	if *serviceTagsMap {
		writeTemplate(templateBody.serviceTagsMap, "servicetagsmap", templateData)
	}

	if *serviceTagsSlice {
		writeTemplate(templateBody.serviceTagsSlice, "servicetagsslice", templateData)
	}

	if *updateTags {
		writeTemplate(templateBody.updateTags, "updatetags", templateData)
	}
}

func writeTemplate(body string, templateName string, td TemplateData) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", filename, err)
	}

	tplate, err := template.New(templateName).Parse(body)
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tplate.Execute(&buffer, td)
	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	contents, err := format.Source(buffer.Bytes())
	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	if _, err := f.Write(contents); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("error closing file (%s): %s", filename, err)
	}
}

func ToSnakeCase(str string) string {
	result := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	result = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(result, "${1}_${2}")
	return strings.ToLower(result)
}
