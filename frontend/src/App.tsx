import { Routes, Route, Navigate } from 'react-router-dom';
import Navigation from './components/Navigation';
import Pipeline from './components/Pipeline';
import Search from './components/Search';
import DomainList from './components/DomainList';

function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      <main className="py-8">
        <Routes>
          <Route path="/" element={<Navigate to="/pipeline" replace />} />
          <Route path="/pipeline" element={<Pipeline />} />
          <Route path="/search" element={<Search />} />
          <Route path="/onapi" element={<DomainList domain="onapi" title="ONAPI Entities" />} />
          <Route path="/scj" element={<DomainList domain="scj" title="SCJ Cases" />} />
          <Route path="/dgii" element={<DomainList domain="dgii" title="DGII Registers" />} />
          <Route path="/pgr" element={<DomainList domain="pgr" title="PGR News" />} />
          <Route path="/docking" element={<DomainList domain="docking" title="Google Docking Results" />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
