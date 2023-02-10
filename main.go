package main

import ("fmt"
	"path/filepath"
	"io/fs"
	"strings"
	"os/exec"
	"os"
)

// enum definitions

type BaseType int16

const (
	STRING BaseType = 1
	BOOLEAN		= 2
	DURATION	= 3
	TIME		= 4
	FLOAT		= 5
	INT64		= 6
	INT32		= 7
	DECIMAL		= 8
	GUID		= 9
)

// struct definitions

type EntityType struct {
	name string
	properties []string
	basetype string
	basetype_flat string
	namespace string
	abstract string
	from_file string
}

type Enum struct {
	name string
	values []string
	namespace string
	from_file string
}

type Complex struct {
	name string
	namespace string
	from_file string
}

type TypeDef struct {
	name string
	basetype string
	namespace string
	from_file string
}

type Collection struct {
	name string
	contained_type TypeDef
	from_file string
}

type Property struct {
	name string
	_type TypeDef
	permissions string
	description string
	long_description string
	from_file string
}

type NavigationProperty struct {
	name string
	_type TypeDef
	permissions string
	auto_expand bool
	expand_references bool
	description string
	long_description string
	from_file string
	contains_target bool
}

// other

var circular_imports = []string{
	"SubProcessors",
	"AllocatedPools",
	"CapacitySources",
	"StorageGroups",
	"Steps",
	"MetricReportDefinition",
	"SubTasks",
	"DataProtectionLinesOfService",
}

// function definitions

func basetype_from_edm(edm_type string) BaseType {
	var e = edm_type
	if e == "Edm.String" { return STRING }
	if e == "Edm.Boolean" { return BOOLEAN }
	if e == "Edm.Decimal" { return DECIMAL }
	if e == "Edm.Int64" { return INT64 }
	if e == "Edm.Int32" { return INT32 }
	if e == "Edm.DateTimeOffset" { return TIME }
	if e == "Edm.Duration" { return DURATION }
	if e == "Edm.GUID" { return GUID }
	return STRING
}

func basetype_to_grpc(b BaseType) (string, []string) {
	var e = []string{}
	switch b {
	case STRING: return "string", e
	case BOOLEAN: return "bool", e
	case DECIMAL: return "double", e
	case INT64: return "int64", e
	case INT32: return "int32", e
	case TIME: return "google.protobuf.Timestamp", []string{"google/protobuf/timestamp.proto"}
	case DURATION: return "google.protobuf.Duration", []string{"google/protobuf/duration.proto"}
	case GUID: return "string", e
	default: return "Can't find type for TODO", e
	}
}

func init_Property() {
	//TODO
}

func init_NavigationProperty() {
	//TODO
}

func find_element_in_scope(element_name string, references []string, this_file string) {
	//TODO
}

func parse_file(filename string, namespaces_to_check []string, element_name_filter string) string {
	//TODO
	return ""
}

func parse_toplevel(filepath string) string {
	fmt.Printf("Parsing %s\n", filepath);
	return parse_file(filepath, []string{}, "");
}

func get_grpc_filename_from_entity(entity string) {
	//TODO
}

func get_grpc_property_type_string(object_type string, this_package string) {
	//TODO
}

func generate_properties_for_entity(typedef string, index_start int, message_name string, package_name string) {
	//TODO
}

func generate_grpc_for_type(typedef string) {
	//TODO
}

func write_fixed_messages() {
	var nav_out = "syntax = \"proto3\";\n\n"
	nav_out += "message NavigationReference {\n"
	nav_out += "    string id = 1;\n }"

	var filepath = "grpc/NavigationReference.proto"
	os.WriteFile(filepath, []byte(nav_out), 0644)
	write_meson_file_for_proto(filepath)
}

func clear_and_make_output_dirs() {

	var info, _ = os.Stat("grpc")

	if info.IsDir() {

		filepath.WalkDir("grpc", func(path string, d fs.DirEntry, err error) error {

			if !d.Type().IsDir() {
				os.Remove(path)
			} else {
				os.RemoveAll(path)
			}

			return nil
		})
	} else { os.Mkdir("grpc", os.ModePerm) }

	var info2, _ = os.Stat("grpc/proto_out")

	if info2.IsDir() {

		filepath.WalkDir("grpc/proto_out", func(path string, d fs.DirEntry, err error) error {

			var info3,_ = os.Stat(path)
			if info3.IsDir() { os.Remove(path) } else { os.RemoveAll(path) }

			return nil
		})
	} else { os.Mkdir("grpc/proto_out", os.ModePerm) }
}

func get_properties_for_service_root(entity string, path string, depth int, collectionlist []string) {
	//TODO
}

func write_service_root(flat_list []string) {
	//TODO
}

func get_cpp_for_type(name string, property_obj string, indent_level int, val_level int, main_object_available bool) {
	//TODO
}

func generate_cpp_for_entitiy(entity string, path string, url_path string, inherit_depth int, collectionlist []string) {
	//TODO
}

func write_cpp_code(flat_list []string) {
	//TODO
}

func write_meson_root_config() {
	//TODO
}

func write_meson_file_for_proto(inputpath string) {
	//TODO
}

func get_lowest_type(this_class string, depth int) {
	//TODO
}

func find_type_for_abstract(class_list []string, abs int) {

	//TODO
}

func instantiate_abstract_classes(class_list []string, this_class string) {

	//TODO
}

func remove_old_schemas(flat_list []string) []string {
	//remove all but the last schema version for a type by
	//loading them into a map with Namespace+name
	var elements = make(map[string]string)
	//for _,item := range flat_list {
		//elements[item.namespace.split(".")[0] + item.name] = item
	//}
	var flat_list_new = []string{}
	for _,v := range elements {
		flat_list_new = append(flat_list_new, v)
	}
	return flat_list_new
}

func main() {

	var gen_protos bool = true
	var gen_cpp bool = true

	if gen_protos {

		var flat_list = []string{}

		fmt.Printf("Reading from csdl")

		var filepaths = []string{}

		filepath.WalkDir("csdl", func(path string, d fs.DirEntry, err error) error {

			if !strings.HasPrefix(d.Name(), "OemAccountService") {
				filepaths = append(filepaths, path+"/"+d.Name())
			}

			//fmt.Printf("file %s\n", d.Name())
			fmt.Printf("file %s\n", path)
			return nil
		});

		for _,filepath := range filepaths {
			flat_list = append(flat_list, parse_toplevel(filepath))
		}

		flat_list = remove_old_schemas(flat_list)

		//TODO: some sorting going on here

		for _,element := range flat_list {
			instantiate_abstract_classes(flat_list, element)
		}

		clear_and_make_output_dirs()

		for _,thistype := range flat_list {
			generate_grpc_for_type(thistype)
		}

		write_fixed_messages()

		write_service_root(flat_list)
		write_cpp_code(flat_list)
		write_meson_root_config()
	}

	if gen_cpp {

		var args = []string{}

		filepath.WalkDir("grpc", func(path string, d fs.DirEntry, err error) error {

			if !strings.HasSuffix(d.Name(), ".proto") { return nil }

			args = []string{"protoc", "--cpp_out", "proto_out", "-I", "grpc", "grpc"+"/"+path}

			return nil
		})

		for _,s := range args { fmt.Printf("%s ", s) }
		fmt.Printf("\n")

		var cmd string = strings.Join(args[:], " ")

		var _, _ = exec.Command(cmd).Output()
	}
}
