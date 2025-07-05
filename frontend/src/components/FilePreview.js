const FilePreview = ({ fileContent, searchResults }) => {
  const renderFileContent = () => {
    if (!searchResults) {
      return <pre>{fileContent}</pre>;
    }

    const lines = fileContent.split("\n");
    const resultsByLine = searchResults.reduce((acc, result) => {
      if (!acc[result.lineNumber]) {
        acc[result.lineNumber] = [];
      }
      acc[result.lineNumber].push(result);
      return acc;
    }, {});

    return (
      <pre>
        {lines.map((line, i) => {
          const lineResults = resultsByLine[i];
          if (!lineResults) {
            return <div key={i}>{line}</div>;
          }

          let lastIndex = 0;
          const parts = [];
          lineResults.forEach((result, j) => {
            parts.push(line.substring(lastIndex, result.startIndex));
            parts.push(
              <mark key={j}>
                {line.substring(result.startIndex, result.endIndex)}
              </mark>,
            );
            lastIndex = result.endIndex;
          });
          parts.push(line.substring(lastIndex));

          return <div key={i}>{parts}</div>;
        })}
      </pre>
    );
  };

  return (
    <div>
      <h2>File Preview</h2>
      {renderFileContent()}
      {searchResults && (
        <div>
          <h3>Search Results</h3>
          <ul>
            {searchResults.map((result, index) => (
              <li key={index}>
                <strong>Line {result.lineNumber}: </strong>
                <span>
                  {result.line.substring(0, result.startIndex)}
                  <mark>
                    {result.line.substring(result.startIndex, result.endIndex)}
                  </mark>
                  {result.line.substring(result.endIndex)}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default FilePreview;
