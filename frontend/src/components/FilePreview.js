const FilePreview = ({ fileContent, searchResults }) => {
  const renderFileContent = () => {
    if (!searchResults) {
      return <pre>{fileContent}</pre>;
    }

    const lines = fileContent.split(/\r?\n/);
    const resultsByLine = searchResults.reduce((acc, result) => {
      if (!acc[result.lineNumber]) {
        acc[result.lineNumber] = [];
      }
      acc[result.lineNumber].push(result);
      return acc;
    }, {});

    return (
      <pre className="whitespace-pre-wrap">
        {lines.map((line, i) => {
          if (line.trim() === '') {
            return <div key={i}>{'\u00A0'}</div>;
          }
          const lineResults = resultsByLine[i];
          if (!lineResults) {
            return <div key={i}>{line}</div>;
          }

          let lastIndex = 0;
          const parts = [];
          lineResults.forEach((result, j) => {
            parts.push(line.substring(lastIndex, result.startIndex));
            parts.push(
              <mark className="bg-yellow-300" key={j}>
                {line.substring(result.startIndex, result.endIndex)}
              </mark>
            );
            lastIndex = result.endIndex;
          });
          parts.push(line.substring(lastIndex));

          return <div key={i}>{parts}</div>;
        })}
      </pre>
    );
  };

  const renderSearchResults = () => (
    <div className="mt-4">
      <h3 className="text-lg font-bold mb-2">Search Results</h3>
      <ul>
        {searchResults.map((result, index) => (
          <li key={index} className="mb-1">
            <strong>Line {result.lineNumber}: </strong>
            <span>
              {result.line.slice(0, result.startIndex)}
              <mark className="bg-yellow-300">
                {result.line.slice(result.startIndex, result.endIndex)}
              </mark>
              {result.line.slice(result.endIndex)}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );

  return (
    <div className="border rounded-lg p-4 mt-4">
      <h2 className="text-xl font-bold mb-2">File Preview</h2>
      <div className="bg-gray-50 p-4 rounded">{renderFileContent()}</div>
      {searchResults?.length > 0 && renderSearchResults()}
    </div>
  );
};

export default FilePreview;
