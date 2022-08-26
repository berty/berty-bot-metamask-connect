import './index.css';
import BaseApp from './pages/App.tsx';
import reportWebVitals from './reportWebVitals';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Sign from "./pages/Sign.tsx";
import SigDisplay from "./pages/defaultcb";

const React = require('react');
const ReactDOM = require('react-dom/client');

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<BaseApp/>}/>
            <Route path="/sign" element={<Sign/>}/>
            <Route path="/default" element={<SigDisplay/>}/>
        </Routes>
    </BrowserRouter>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
