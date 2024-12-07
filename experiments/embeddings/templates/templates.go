package templates

import "html/template"

// Templates contains all the HTML templates used in the application
var Templates = struct {
	Base          *template.Template
	SearchResults *template.Template
}{
	Base:          template.Must(template.New("base").Parse(baseHTML)),
	SearchResults: template.Must(template.New("results").Parse(searchResultsHTML)),
}

const baseHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Bleve Search</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 p-8">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-3xl font-bold mb-8">Bleve Search Demo</h1>
        
        <div class="bg-white rounded-lg shadow p-6 mb-8">
            <h2 class="text-xl font-semibold mb-4">Index Mappings</h2>
            <pre class="bg-gray-50 p-4 rounded overflow-auto max-h-96"><code>{{.Mappings}}</code></pre>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
            <h2 class="text-xl font-semibold mb-4">Search</h2>
            <form hx-post="/search" hx-target="#results" class="mb-6">
                <div class="flex gap-4">
                    <input type="text" name="query" placeholder="Enter your search query" 
                           class="flex-1 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <button type="submit" 
                            class="px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">
                        Search
                    </button>
                </div>
            </form>
            <div id="results"></div>
        </div>
    </div>
</body>
</html>`

const searchResultsHTML = `
{{if .Error}}
    <div class="text-red-500 mb-4">{{.Error}}</div>
{{else}}
    <div class="space-y-4">
        {{range .Hits}}
            <div class="border rounded p-4 hover:bg-gray-50">
                <div class="font-semibold mb-2">Document ID: {{.ID}}</div>
                <div class="text-gray-600">Score: {{printf "%.4f" .Score}}</div>
                {{if .Fragments}}
                    <div class="mt-2 text-sm">
                        {{range $field, $fragments := .Fragments}}
                            {{range $fragments}}
                                <div class="mt-1">... {{.}} ...</div>
                            {{end}}
                        {{end}}
                    </div>
                {{end}}
            </div>
        {{else}}
            <div class="text-gray-500">No results found</div>
        {{end}}
    </div>
{{end}}`
