/* import fs from "fs";

// Load package.json
const packageJson = JSON.parse(fs.readFileSync("package.json", "utf-8"));
const dependencies = packageJson.dependencies || {};
const devDependencies = packageJson.devDependencies || {};

// Combine dependencies
const allDependencies = { ...dependencies, ...devDependencies };

// Function to generate clusters and edges
function generateDot(dependencies) {
  let dot = `
digraph G {
    // Set overall graph attributes
    rankdir=TB;
    layout=dot;
    splines=polyline;
    overlap=false;
    concentrate=true;
    sep=0.3;
    nodesep=0.2;
    ranksep=1;
    node [shape=rectangle, style=filled, fontname="Helvetica", fontsize=10];

    // Define clusters for better organization

    subgraph cluster_project {
        label = "Project Packages";
        color = lightblue;
        style=filled;
        // List your frontend project packages
        "frontend";
        "frontend/components";
        "frontend/utils";
        "frontend/store";
        "frontend/routes";
    }

    subgraph cluster_stdlib {
        label = "Standard Library";
        color = lightgrey;
        style=filled;
        rank=same;  // Arrange nodes horizontally
        // Standard library nodes (for frontend, typically browser APIs)
        "fetch";
        "URL";
        "console";
    }

    subgraph cluster_third_party {
        label = "Third-Party Packages";
        color = beige;
        style=filled;
`;

  // Add third-party packages
  for (const pkg of Object.keys(dependencies)) {
    dot += `           "${pkg}";\n`;
  }

  dot += `       }\n\n`;

  // Define project dependencies (manual or based on your analysis)
  // Example:
  dot += `
    // Edges representing dependencies
    "frontend" -> "frontend/components";
    "frontend/components" -> "svelte";
    "frontend/utils" -> "lodash";
    "frontend/routes" -> "svelte-routing";
    "frontend" -> "axios";
    "frontend" -> "moment";
    "frontend" -> "fetch";
    "frontend" -> "URL";
    "frontend" -> "console";
`;

  dot += "}";
  return dot;
}

const dotContent = generateDot(allDependencies);

fs.writeFileSync("frontend-dependency-graph.dot", dotContent, "utf-8");
console.log("Frontend dependency graph generated: frontend-dependency-graph.dot");
 */
