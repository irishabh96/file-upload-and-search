import React from "react";

const FileList = ({ files, onSelectFile }) => {
  return (
    <div>
      <h2>Files</h2>
      {files && (
        <ul>
          {files?.map((file) => (
            <li key={file.ID} onClick={() => onSelectFile(file)}>
              {file.Name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default FileList;
