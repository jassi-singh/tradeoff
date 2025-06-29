const OverlayChildren = () => {
  return (
    <>
      <h2 className="text-2xl font-bold text-white">Enter the Arena</h2>
      <input
        type="text"
        name="username"
        placeholder="Enter your name"
        className="p-2 rounded-md bg-gray-700 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        type="submit"
        className="bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 transition-colors"
      >
        Join Game
      </button>
    </>
  );
};

export default OverlayChildren;
