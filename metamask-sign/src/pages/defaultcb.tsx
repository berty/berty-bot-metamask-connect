import React, {useEffect} from "react";
import "../index.css";
import Navbar from "../components/navbar";
import Display from "../components/display";

function SigDisplay() {
    const [sig, setSig] = React.useState<string|null>(null);

    useEffect(() => {
        const sigParam = new URLSearchParams(window.location.search).get('sig');
        if (sigParam === null) {
            console.error("no sig found");
            return;
        }
        setSig(sigParam);
    }, []);

    return (
        <div className="App">
            <Navbar/>
            <Display fieldName={'Signature'} content={sig !== null ? sig : "no sig"}/>
        </div>
    );
}

export default SigDisplay;
