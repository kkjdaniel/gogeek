package request

import (
	"strings"
	"testing"
)

func TestFixMalformedXMLMultilineCDATA(t *testing.T) {
	input := []byte("<item><desc><![CDATA[line one & <tag>\nline two & more]]></desc></item>")
	out := string(fixMalformedXML(input))
	if !strings.Contains(out, "<![CDATA[line one & <tag>\nline two & more]]>") {
		t.Errorf("multi-line CDATA was mangled: %s", out)
	}
}
