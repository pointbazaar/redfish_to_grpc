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
	properties []Property
	basetype *EntityType
	basetype_flat string
	namespace string
	abstract string
	from_file string
	contained_type *EntityType
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
	contained_type *Collection
	from_file string
}

type Property struct {
	name string
	_type EntityType
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

var folders_added_to_grpc = []string{}

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

func init_Property(p Property, name string, thistype EntityType, permissions string, description string, long_description string, from_file string) {

	p.name = name;
	p._type = thistype;
	p.permissions = permissions;
	p.description = description;
	p.long_description = long_description;
	p.from_file = from_file;
}

func init_NavigationProperty(p NavigationProperty, name string, thistype TypeDef, permissions string, auto_expand bool, expand_references bool, description string, long_description string, from_file string, contains_target bool) {

	p.name = name;
	p._type = thistype;
	p.permissions = permissions;
	p.description = description;
	p.long_description = long_description;
	p.from_file = from_file;
	p.auto_expand = auto_expand;
	p.expand_references = expand_references;
	p.contains_target = contains_target;
}

func find_element_in_scope(element_name string, references []string, this_file string) Collection {

	var contained_element Collection = Collection{}

	if strings.HasPrefix(element_name, "Collection(") && strings.HasSuffix(element_name, ")") {

		contained_element = find_element_in_scope(element_name[11:], references, this_file)
		return Collection{name: element_name, contained_type: &contained_element, from_file: this_file}
	}

	//var edmtype = basetype_from_edm(element_name) //TODO

	//if edmtype != nil { return edmtype }

	var elements Collection = Collection{}

	var namespaces = []string{} //TODO
	//for _,(reference_uri, namespaces) := references { //TODO
	for i,_ := range references {

		fmt.Printf("%d", i); //DUMMY
		var reference_uri = "" //TODO

		//var uri = urlparse(reference_uri) //TODO: urlparse function
		var uri = struct {path string ""}{}

		var path = "cdsl" + "/" + filepath.Base(uri.path)

		_,err :=  os.Stat(path)
		if err != nil {

			fmt.Printf("File doesn't exist, downloading from %s\n", reference_uri)
			//TODO: find some requests library
			//var r = requests.get(reference_uri)
			//r.raise_for_status()

			//_ = os.WriteFile(path, []byte(r.content))
		}

		elements = parse_file(path, namespaces, element_name)

		//var l = len(elements) //TODO
		var l = 0

		if l == 0 { continue }
		if l  >  1 {

			fmt.Printf("Found %d %s elemens with referencelist %s\n", l, element_name, "TODO") //TODO
			continue
		}

		//return elements[0]//TODO
		return elements
	}

	//finish by searching the file we're in now
	elements = parse_file(this_file, namespaces, element_name)

	//var l1 = len(elements)//
	var l1 = 1;

	if l1 != 1 { return Collection{} }

	//return elements[0] //TODO
	return elements

	fmt.Printf("Unable to find %s\n", element_name)

	return Collection{}
}

func parse_file(filename string, namespaces_to_check []string, element_name_filter string) Collection {
	//TODO
	return Collection{}
}

func parse_toplevel(filepath string) Collection {
	fmt.Printf("Parsing %s\n", filepath);
	return parse_file(filepath, []string{}, "");
}

func get_grpc_filename_from_entity(entity EntityType) string{

	var new_filename = filepath.Base(entity.from_file)
	var filepath = new_filename + "/" + entity.name + ".proto"
	return filepath
}

func get_grpc_property_type_string(object_type EntityType, this_package string) (string, []string) {

	var required_imports = []string{}

	//TODO
	var filename = get_grpc_filename_from_entity(object_type)

	required_imports = append(required_imports, filename)

	if this_package == strings.Split(object_type.namespace, ".")[0] {
		return object_type.name, required_imports
	}

	return "." + strings.Split(object_type.namespace, ".")[0] + "." + object_type.name, required_imports
}

func generate_properties_for_entity(typedef EntityType, index_start int, message_name string, package_name string) (string, []string, int) {

	var grpc_out = ""
	var required_imports = []string{}
	var property_index = index_start

	if true { //TODO

		if typedef.basetype != nil {
			var text,includes,p1 = generate_properties_for_entity(
				*typedef.basetype, index_start, message_name, package_name)
			property_index = p1
			grpc_out += text
			required_imports = append(required_imports, includes...)
		}

		if len(typedef.properties) != 0 {
			if property_index != 1 {
				grpc_out += "\n"
			}
			grpc_out += fmt.Sprintf("   // from %s.%s\n", typedef.namespace, typedef.name)
		}

		for _,property_obj := range typedef.properties {

			//if isinstance ... //TODO
			if true {
				var grpc_type = "NavigationReference"
				//if isinstance ... //TODO
				grpc_out += fmt.Sprintf("    %s %s = %s;\n", grpc_type, property_obj.name, property_index)
				required_imports = append(required_imports, "NavigationReference.proto")
			} else {
				var text, imports = get_grpc_property_type_string(property_obj._type, package_name)

				grpc_out += fmt.Sprintf("    %s %s = %s;\n", text, property_obj.name, property_index)

				required_imports = append(required_imports, imports...)
			}
			property_index += 1
		}
	}
	return grpc_out, required_imports, property_index
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

func get_properties_for_service_root(entity *EntityType, path string, depth int, collectionlist []string) (string, []string, []string) {

	var body = "";
	var header = []string{};
	var messages = []string{};

	//TODO(ed) What do collections look like?

	//if isinstance(entity, Collection)
		collectionlist = []string{};
		collectionlist = append(collectionlist, path);
		return get_properties_for_service_root(entity.contained_type, path, 0, collectionlist);

	//if isinstance(entity, TypeDef)
		//Note, because typedefs aren't a real property, they don't increase depth
		return get_properties_for_service_root(entity.basetype, path, 0, collectionlist);

	//if isinstance(entity, EntityType)
		if path == "" { path = "ServiceRoot";  }

		// only generate service definition for root level objects
		if depth == 0 {

			var s = fmt.Sprintf("\nmessage Get_%s_FilterSpec{{\n   string expand = 1;\n    repeated string filter = 2;\n", path);
			messages = append(messages, s);

			for element_index, _ := range collectionlist {

				var arr = strings.Split(path, "_");
				var s = fmt.Sprintf("    NavigationReference %sId = %s;\n", arr[len(arr)-1], element_index+3);
				messages = append(messages, s);

				//TODO(ed) Figure out how to rename routes with multiple ids
				break;
			}

			messages = append(messages, "}\n");

			var ns1 = strings.Split(entity.namespace, ".")[0];
			body += (fmt.Sprintf("    rpc Get_%s(Get_%s_FilterSpec) returns (%s.%s) {{}};\n", path, path, ns1, entity.name));

			var filename = get_grpc_filename_from_entity(*entity);

			header = append(header, fmt.Sprintf("import \"%s\";\n", filename));
		}

		if entity.basetype != nil {

			var this_body, this_header, this_messages = get_properties_for_service_root(entity.basetype, path, depth+1, collectionlist);
			body += this_body;
			header = append(header, this_header...);

			messages = append(messages, strings.Join(this_messages, ""));
		}

		for _,property_obj := range entity.properties {

			//if !isinstance(property_obj, NavigationProperty {continue;}
			var new_path = path + "_" + property_obj.name;

			if strings.HasPrefix(new_path, "_") { new_path = new_path[1:]; }

			body += fmt.Sprintf("\n    //from %s\n", entity.namespace + "." + entity.name);

			var this_body, this_header, this_messages = get_properties_for_service_root(&property_obj._type, new_path, 0, collectionlist);

			body += this_body;
			header = append(header, this_header...);
			messages = append(messages, strings.Join(this_messages, ""));
		}

	return body, header, messages;
}

func write_service_root(flat_list []string) {

	//service_root = [x for x in flat_list if x.name == "ServiceRoot"] //TODO
	var service_root = []*EntityType{}

	if len(service_root) != 1 {
		fmt.Printf("Unable to find unique service root\n")
		//TODO: raise exception
	}

	var proto_filename = "grpc/entry.proto"

	var f,_ = os.Create(proto_filename)
	defer f.Close()

	f.WriteString("syntax = \"proto3\";\n\n");
	f.WriteString("package redfish_v1;\n\n")

	var body, header, messages = get_properties_for_service_root(service_root[0], "", 0, []string{})

	header = append(header, "import \"NavigationReference.proto\";\n")
	//deduplicate headers
	//header = sorted(set(header), key=str.casefold) //TODO

	for _,header_element := range header {
		f.WriteString(header_element)
	}

	for _,m := range messages { f.WriteString(m) }

	f.WriteString("service Redfish_v1{\n")
	f.WriteString(body)
	f.WriteString("}")

	write_meson_file_for_proto(proto_filename)
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

	var filepath = "grpc/meson.build"
	var f,_ = os.Create(filepath)
	defer f.Close()

	f.WriteString("protobuf_generated = []\n")

	//for filename in sorted(folders_added_to_grpc) //TODO
		//TODO

	f.WriteString("protobuf_generated += grpc_gen.process(\n")
	f.WriteString("'entry.proto',\n")
	f.WriteString("preserve_path_from : meson.current_source_dir()")
	f.WriteString(")\n")
}

func write_meson_file_for_proto(inputpath string) {

	//var dirname = filepath.Dir(inputpath)
	//var basename = filepath.Base(inputpath)
	//var meson_filename = dirname + "meson.build"
	var dirpath,_ = filepath.Rel(inputpath, "grpc")

	if dirpath != "." {
		folders_added_to_grpc = append(folders_added_to_grpc, dirpath)
	}
}

func get_lowest_type(this_class EntityType, depth int) EntityType {
	//TODO
	//return EntityType{}
	return this_class
}

func find_type_for_abstract(class_list []EntityType, abs EntityType) EntityType {

	for _,element := range class_list {

		var lt EntityType = get_lowest_type(element, 0)
		if lt.name == abs.name && lt.from_file == abs.from_file { return element }
	}

	return abs
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

			//flat_list = append(flat_list, parse_toplevel(filepath)) //TODO
			flat_list = append(flat_list, filepath)
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
