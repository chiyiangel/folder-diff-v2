package report

import (
	"html/template"
	"os"

	"folder-diff-v2/internal/compare"
)

type Generator struct {
	template *template.Template
}

func NewGenerator() *Generator {
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))
	return &Generator{template: tmpl} // Fixed syntax
}

func (g *Generator) GenerateHTML(result *compare.ComparisonResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.template.Execute(file, result)
}

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Folder Comparison Report</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .tree { margin-left: 20px; }
        .identical { color: green; }
        .modified { color: red; }
        .new { color: blue; }
        .deleted { color: gray; }
        .folder { font-weight: bold; cursor: pointer; }
        .legend { margin: 20px 0; padding: 10px; background: #f5f5f5; border-radius: 4px; }
        .content { display: flex; gap: 20px; }
        .source, .target { flex: 1; }
        .file-details { 
            display: none; 
            position: fixed; 
            top: 50%; 
            left: 50%; 
            transform: translate(-50%, -50%); 
            background: white;
            padding: 20px; 
            box-shadow: 0 0 10px rgba(0,0,0,0.3);
            border-radius: 4px; 
            z-index: 1000;
        }
        .close-button { float: right; cursor: pointer; }
        .tree ul { list-style: none; padding-left: 20px; }
        .tree li { margin: 5px 0; }
        .icon { margin-right: 5px; width: 16px; }
        .collapsed ul { display: none; }
    </style>
</head>
<body>
    <h1><i class="fas fa-folder-open"></i> Folder Comparison Report</h1>
    <div class="legend">
        <div><i class="fas fa-check identical"></i> Identical</div>
        <div><i class="fas fa-exclamation-triangle modified"></i> Modified</div>
        <div><i class="fas fa-plus new"></i> New</div>
        <div><i class="fas fa-minus deleted"></i> Deleted</div>
    </div>
    <div class="content">
        <div class="source">
            <h2><i class="fas fa-folder"></i> Source: {{.SourceRoot}}</h2>
            <div class="tree">
                {{template "fileTree" .SourceFiles}}
            </div>
        </div>
        <div class="target">
            <h2><i class="fas fa-folder"></i> Target: {{.TargetRoot}}</h2>
            <div class="tree">
                {{template "fileTree" .TargetFiles}}
            </div>
        </div>
    </div>
    <div id="fileDetails" class="file-details">
        <i class="fas fa-times close-button" onclick="closeDetails()"></i>
        <h3>File Details</h3>
        <div id="fileContent"></div>
    </div>
    <script>
        function toggleFolder(element) {
            const li = element.closest('li');
            li.classList.toggle('collapsed');
        }

        function showFileDetails(path, status, hash) {
            const details = document.getElementById('fileDetails');
            const content = document.getElementById('fileContent');
            content.innerHTML = 
                '<p><strong>Path:</strong> ' + path + '</p>' +
                '<p><strong>Status:</strong> <span class="' + status + '">' + status + '</span></p>' +
                (hash ? '<p><strong>Hash:</strong> ' + hash + '</p>' : '');
            details.style.display = 'block';
        }

        function closeDetails() {
            document.getElementById('fileDetails').style.display = 'none';
        }

        document.addEventListener('DOMContentLoaded', function() {
            document.querySelectorAll('.folder').forEach(function(folder) {
                folder.classList.add('collapsed');
            });
        });
    </script>
</body>
</html>

{{define "fileTree"}}
<ul>
    {{range .}}
    <li class="{{if .IsDir}}folder{{else}}file{{end}}">
        {{if .IsDir}}
            <span onclick="toggleFolder(this)">
                <i class="fas fa-folder icon"></i>{{.RelPath}}
            </span>
            <ul>
                {{template "fileTree" .Children}}
            </ul>
        {{else}}
            <span class="{{.Status}}" onclick="showFileDetails('{{.RelPath}}', '{{.Status}}', '{{.Hash}}')">
                <i class="fas fa-file icon"></i>{{.RelPath}}
            </span>
        {{end}}
    </li>
    {{end}}
</ul>
{{end}}`
