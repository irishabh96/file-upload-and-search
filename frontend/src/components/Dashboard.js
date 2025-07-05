import React, { useState, useEffect } from "react";
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

  const fetchFiles = async () => {
    try {
      const response = await axios.get("/api/files");
      setFiles(response.data);
    } catch (error) {
      if (error.response && error.response.status === 401) {
        navigate("/login");
      }
      console.error("Error fetching files:", error);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

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
    if (!selectedFile) return;

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
    <div>
      <h1>Dashboard</h1>
      <Upload onUploadSuccess={fetchFiles} />
      <FileList files={files} onSelectFile={handleFileSelect} />
      {selectedFile && <Search onSearch={handleSearch} />}
      <FilePreview fileContent={fileContent} searchResults={searchResults} />
    </div>
  );
};

export default Dashboard;
