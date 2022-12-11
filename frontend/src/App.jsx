import { BrowserRouter as Router, Route, } from 'react-router-dom'
import Register from './pages/Register'
import Login from './pages/Login'
import Session from './pages/Session'
import './App.css'

function App() {
  return (
    <Router>
      <Route exact path="/" component={Session} />
      <Route exact path="/user/register" component={Register} />
      <Route exact path="/user/login" component={Login} />
      <Route exact path="/session" component={Session} />
    </Router>
  )
}

export default App
