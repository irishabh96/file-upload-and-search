import React, { useState } from "react";
import axios from "axios";

const Upload = ({ onUploadSuccess }) => {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);

    try {
      await axios.post("/api/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      alert("File uploaded successfully");
      onUploadSuccess();
    } catch (error) {
      console.error("Upload failed", error);
    }
  };

  return (
    <div>
      <h2>Upload</h2>
      <form onSubmit={handleUpload}>
        <input required type="file" onChange={handleFileChange} />
        <button type="submit">Upload</button>
      </form>
    </div>
  );
};

export default Upload;
