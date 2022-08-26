import "../index.css";
import React from "react";
import Navbar from "../components/navbar";
import FormSign from "../components/form";

function BaseApp() {
    return (
        <div className="App">
            <Navbar/>
            <FormSign/>
        </div>
    );
}

export default BaseApp;
