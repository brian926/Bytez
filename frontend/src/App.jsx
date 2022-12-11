import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
// import Register from './pages/Register'
// import Login from './pages/Login'
import Session from './pages/Session2'
import './App.css'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Session />} />
        <Route path="/session" element={<Session />} />
      </Routes>
    </Router>
  )
}

export default App
