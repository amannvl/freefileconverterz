import { Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import Convert from './pages/Convert';
import Terms from './pages/Terms';
import Privacy from './pages/Privacy';
import About from './pages/About';
import Help from './pages/Help';
import Contact from './pages/Contact';
import CookiePolicy from './pages/CookiePolicy';
import GDPR from './pages/GDPR';
import NotFound from './pages/NotFound';

function App() {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <Header />
      <main className="flex-grow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="convert/:conversionType" element={<Convert />} />
            <Route path="terms" element={<Terms />} />
            <Route path="privacy" element={<Privacy />} />
            <Route path="about" element={<About />} />
            <Route path="help" element={<Help />} />
            <Route path="contact" element={<Contact />} />
            <Route path="cookie-policy" element={<CookiePolicy />} />
            <Route path="gdpr" element={<GDPR />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </div>
      </main>
      <Footer />
    </div>
  );
}

export default App;
