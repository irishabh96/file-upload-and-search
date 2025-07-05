import React from "react";

const FileList = ({ files, onSelectFile }) => {
  return (
    <div className="border rounded-lg p-4">
      <h2 className="text-xl font-bold mb-2">Files</h2>
      {files && (
        <ul>
          {files?.map((file) => (
            <li
              key={file.ID}
              onClick={() => onSelectFile(file)}
              className="cursor-pointer p-2 hover:bg-gray-100 rounded"
            >
              {file.Name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default FileList;
