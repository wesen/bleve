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
            <h2 class="text-xl font-semibold mb-4">Embeddings Information</h2>
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <h3 class="text-lg font-medium mb-2">Model Information</h3>
                    <p><span class="font-semibold">Model:</span> {{.EmbeddingsInfo.model}}</p>
                    <p><span class="font-semibold">Dimensions:</span> {{.EmbeddingsInfo.dimensions}}</p>
                </div>
                <div>
                    <h3 class="text-lg font-medium mb-2">Cache Information</h3>
                    <p><span class="font-semibold">Type:</span> {{.EmbeddingsInfo.cache.type}}</p>
                    {{if eq .EmbeddingsInfo.cache.type "memory"}}
                    <p><span class="font-semibold">Size:</span> {{.EmbeddingsInfo.cache.size}} / {{.EmbeddingsInfo.cache.max_size}}</p>
                    {{end}}
                    <div class="mt-2">
                        <button hx-post="/cache/clear" 
                                hx-target="#cache-status" 
                                class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500">
                            Clear Cache
                        </button>
                        <button hx-get="/cache/stats" 
                                hx-target="#cache-status" 
                                class="ml-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">
                            Refresh Stats
                        </button>
                    </div>
                    <div id="cache-status" class="mt-2 text-sm"></div>
                </div>
            </div>
        </div>
        
        <div class="bg-white rounded-lg shadow p-6 mb-8">
            <h2 class="text-xl font-semibold mb-4">Index Mappings</h2>
            <pre class="bg-gray-50 p-4 rounded overflow-auto max-h-96"><code>{{.Mappings}}</code></pre>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
            <h2 class="text-xl font-semibold mb-4">Search</h2>
            <div class="flex gap-4 mb-4">
                <div class="flex-1">
                    <button class="w-full py-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            hx-get="/documents"
                            hx-target="#results">
                        List All Documents
                    </button>
                </div>
                <div class="flex-1">
                    <button class="w-full py-2 bg-green-500 text-white rounded hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
                            onclick="toggleEmbeddingForm()">
                        Generate Embedding
                    </button>
                </div>
            </div>
            
            <div id="embedding-form" class="hidden mb-6 p-4 border border-gray-200 rounded bg-gray-50">
                <h3 class="text-lg font-medium mb-2">Generate Embedding</h3>
                <form hx-post="/embeddings" hx-target="#embedding-result">
                    <div class="mb-4">
                        <label class="block text-gray-700 mb-2">Text:</label>
                        <textarea name="text" class="w-full px-3 py-2 border rounded" rows="3"></textarea>
                    </div>
                    <button type="submit" class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600">
                        Generate
                    </button>
                </form>
                <div id="embedding-result" class="mt-4"></div>
            </div>
            
            <form hx-post="/search" hx-target="#results" class="mb-6">
                <h3 class="text-lg font-medium mb-2">YAML Query:</h3>
                <div class="mb-4">
                    <textarea name="yaml" class="w-full px-4 py-2 border rounded font-mono" rows="10" placeholder="Enter your YAML query..."></textarea>
                </div>
                <div class="flex justify-between">
                    <button type="submit" 
                            class="px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">
                        Search
                    </button>
                    <button type="button" onclick="loadExampleQuery()"
                            class="px-6 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500">
                        Load Example
                    </button>
                </div>
            </form>
            <div id="results"></div>
        </div>
    </div>
    
    <script>
        function toggleEmbeddingForm() {
            const form = document.getElementById('embedding-form');
            form.classList.toggle('hidden');
        }
        
        function loadExampleQuery() {
            const textarea = document.querySelector('textarea[name="yaml"]');
            textarea.value = 'query:\n  vector:\n    field: vector\n    text: "What is the meaning of life?"\n    model: all-minilm\n    k: 10\n    boost: 1.0\noptions:\n  size: 10\n  highlight:\n    fields: [content]';
        }
    </script>
</body>
</html>`

const searchResultsHTML = `
{{if .Error}}
    <div class="text-red-500 mb-4">{{.Error}}</div>
{{else}}
    <div class="mt-4">
        <div class="mb-2 text-gray-700">Found {{.Total}} results in {{.Took}}</div>
        <div class="space-y-4">
            {{range .Hits}}
                <div class="border rounded p-4 hover:bg-gray-50">
                    <div class="font-semibold mb-2">Document ID: {{.ID}}</div>
                    <div class="text-gray-600">Score: {{printf "%.4f" .Score}}</div>
                    {{if .Fields}}
                        <div class="mt-2 text-sm">
                            {{range $field, $value := .Fields}}
                                {{if ne $field "vector"}}
                                    <div><span class="font-semibold">{{$field}}:</span> {{$value}}</div>
                                {{else}}
                                    <div><span class="font-semibold">{{$field}}:</span> [Vector with {{len $value}} dimensions]</div>
                                {{end}}
                            {{end}}
                        </div>
                    {{end}}
                    {{if .Fragments}}
                        <div class="mt-2 text-sm bg-yellow-50 p-2 rounded">
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
    </div>
{{end}}`
