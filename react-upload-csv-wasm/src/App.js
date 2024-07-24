import React, { useState, useEffect } from 'react';
import FileUploader from './FileUploader';
import { JsonView, allExpanded, darkStyles } from 'react-json-view-lite';
import 'react-json-view-lite/dist/index.css';
import './App.css';

function App() {
  const [csvData, setCsvData] = useState(null);

  useEffect(() => {
    const loadWASM = async () => {
      if (typeof Go === "undefined") {
        console.error('No se encontró el objeto Go. Asegúrate de que wasm_exec.js está incluido en tu HTML.');
        return;
      }

      const go = new window.Go();
      const wasmModule = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
      go.run(wasmModule.instance);
    };
    loadWASM();
  }, []);

  const handleFileLoaded = (data) => {
    const csvString = data.map(row => row.join(',')).join('\n');
    const result = window.processCSV(csvString);
    const parsedResult = JSON.parse(result);
    setCsvData(parsedResult);
  };

  return (
    <div className="App">
      <h1>CSV File Uploader</h1>
      <div className="container">
        <div className="uploader">
          <FileUploader onFileLoaded={handleFileLoaded} />
        </div>
        <div className="json-output">
          {csvData && (
            <div>
              <h2>JSON Data</h2>
              <JsonView data={csvData} shouldExpandNode={allExpanded} style={darkStyles} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
