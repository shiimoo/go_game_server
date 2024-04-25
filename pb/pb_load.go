package pb

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

func LoadProto(dir string) error {
	pbList, err := getProtoFileList(dir)
	if err != nil {
		return err
	}

	// 解析器
	Parser := protoparse.Parser{}
	descs, err := Parser.ParseFiles(pbList...)
	if err != nil {
		return err
	}
	for _, pbfile := range descs {
		for _, msgType := range pbfile.GetMessageTypes() {
			getMsgSet().RegisterMsg(pbfile.GetPackage(), msgType)
		}
	}

	return nil
}

func getProtoFileList(dir string) ([]string, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	pbList := make([]string, 0)
	for _, fs := range list {
		if fs.IsDir() {
			// todo 包递归
			continue
		}
		fName := strings.TrimSpace(fs.Name())
		if !strings.HasSuffix(fName, ".proto") {
			continue
		}
		pbList = append(pbList, path.Join(dir, fName))
	}
	return pbList, nil
}

func main1() {

	Filename := "./proto/msg.proto"

	Parser := protoparse.Parser{}
	//加载并解析 proto文件,得到一组 FileDescriptor
	descs, err := Parser.ParseFiles(Filename)
	if err != nil {
		fmt.Printf("ParseFiles err=%v", err)
		return
	}
	desc := descs[0]
	msg := desc.FindMessage("test.AddFriendReq")
	dmsg := dynamic.NewMessage(msg)
	dmsg.SetFieldByName("keyword", "shimo")

	fmt.Println(dmsg.Marshal())
}
