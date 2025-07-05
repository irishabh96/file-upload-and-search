import React, { useState } from "react";

const Search = ({ onSearch }) => {
  const [query, setQuery] = useState("");

  const handleSearch = (e) => {
    e.preventDefault();
    onSearch(query);
  };

  return (
    <div className="mb-4">
      <h2 className="text-xl font-bold mb-2">Search</h2>
      <form onSubmit={handleSearch} className="flex items-center">
        <input
          type="text"
          placeholder="Query"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="border rounded-l-lg p-2 flex-grow"
        />
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-r-lg"
        >
          Search
        </button>
      </form>
    </div>
  );
};

export default Search;
