import './App.css';
import { ToolListDemo } from './components/ToolListDemo';

function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto">
        <header className="bg-white shadow-sm mb-8">
          <div className="px-6 py-4">
            <div className="uppercase tracking-wide text-sm text-indigo-500 font-semibold">
              Tool Tracker
            </div>
            <h1 className="mt-1 text-2xl font-medium text-black">Inventory Management</h1>
          </div>
        </header>

        <main>
          <ToolListDemo />
        </main>
      </div>
    </div>
  );
}

export default App;
