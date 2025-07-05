import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Upload from "./Upload";
import FileList from "./FileList";
import Search from "./Search";
import FilePreview from "./FilePreview";

const Dashboard = () => {
  const [files, setFiles] = useState([]);
  const [selectedFile, setSelectedFile] = useState(null);
  const [searchResults, setSearchResults] = useState([]);
  const [fileContent, setFileContent] = useState("");
  const navigate = useNavigate();
  const fetchFiles = useCallback(async () => {
    try {
      const response = await axios.get("/api/files");
      setFiles(response.data);
    } catch (error) {
      if (error.response && error.response.status === 401) {
        navigate("/login");
      }
      console.error("Error fetching files:", error);
    }
  }, [navigate]);

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  const handleFileSelect = async (file) => {
    setSelectedFile(file);
    setSearchResults([]);
    try {
      const response = await axios.get(`/api/file?fileId=${file.ID}`);
      setFileContent(response.data);
    } catch (error) {
      console.error("Error fetching file content:", error);
    }
  };

  const handleSearch = async (term) => {
    if (!term) {
      setSearchResults([]);
      return;
    }

    try {
      const response = await axios.get(
        `/api/search?q=${term}&fileId=${selectedFile.ID}`,
      );
      setSearchResults(response.data);
    } catch (error) {
      console.error("Error searching:", error);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-4">Dashboard</h1>
      <Upload onUploadSuccess={fetchFiles} />
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
        <div className="col-span-1">
          <FileList files={files} onSelectFile={handleFileSelect} />
        </div>
        <div className="col-span-2">
          {selectedFile && <Search onSearch={handleSearch} />}
          <FilePreview
            fileContent={fileContent}
            searchResults={searchResults}
          />
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
