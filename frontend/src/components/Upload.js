import React, { useState } from "react";
import axios from "axios";

const Upload = ({ onUploadSuccess }) => {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async (e) => {
    e.preventDefault();
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
      setFile(null);
    } catch (error) {
      console.error("Upload failed", error);
    }
  };

  return (
    <div className="mb-4">
      <h2 className="text-xl font-bold mb-2">Upload</h2>
      <form onSubmit={handleUpload} className="flex items-center">
        <input
          required
          type="file"
          onChange={handleFileChange}
          className="border rounded-l-lg p-2 flex-grow"
        />
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-r-lg"
        >
          Upload
        </button>
      </form>
    </div>
  );
};

export default Upload;
