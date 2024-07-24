import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import Papa from 'papaparse';
import './FileUploader.css'; // Importa el archivo CSS

const FileUploader = ({ onFileLoaded }) => {
  const onDrop = useCallback((acceptedFiles) => {
    const file = acceptedFiles[0];
    const reader = new FileReader();
    reader.onload = () => {
      const text = reader.result;
      Papa.parse(text, {
        complete: (results) => {
          onFileLoaded(results.data);
        }
      });
    };
    reader.readAsText(file);
  }, [onFileLoaded]);

  const { getRootProps, getInputProps } = useDropzone({ onDrop });

  return (
    <div {...getRootProps({ className: 'dropzone' })}>
      <input {...getInputProps()} />
      <p>Drag & drop a CSV file here, or click to select one</p>
    </div>
  );
};

export default FileUploader;
