import './App.css';

function App() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden">
        <div className="p-8">
          <div className="uppercase tracking-wide text-sm text-indigo-500 font-semibold">
            Tool Tracker
          </div>
          <h1 className="mt-2 text-xl font-medium text-black">Frontend Setup Complete!</h1>
          <p className="mt-2 text-gray-500">
            React + Vite + TypeScript + Tailwind CSS v4.1 is ready to go.
          </p>
          <div className="mt-6">
            <button className="bg-indigo-500 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded transition duration-150">
              Get Started
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
