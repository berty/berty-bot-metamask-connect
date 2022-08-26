import {SetStateAction, useState} from "react";
import {sign, hash256} from "../utils/utils";

const React = require('react');

export default function FormSign(): JSX.Element {
    const [message, setMessage] = useState('');
    const [hash, setHash] = useState('');
    const [signature, setSignature] = useState('');

    function handleHash(): void {
        setHash(hash256(message));
    }

    async function handleSign(): Promise<void> {
        const signature = await sign(hash);
        setSignature(signature);
    }

    const handleChangeMessage = (event: { target: { value: SetStateAction<string>; }; }) => {
        setMessage(event.target.value);
    };

    const handleChangeHash = (event: { target: { value: SetStateAction<string>; }; }) => {
        setMessage(event.target.value);
    };

    return (
        <span className="inline-block align-middle ...">
            <form className="w-full max-w-sm" onSubmit={handleHash}>
            <div className="flex items-center border-b border-teal-500 py-2">
                <input
                    className="appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none"
                    type="text" placeholder="to hash" aria-label="Full name" value={message} onChange={handleChangeMessage}/>
                <button
                    className="flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded"
                    type="button" onClick={handleHash}>
                    sha256
                </button>
            </div>
        </form>
        <form className="w-full max-w-sm" onSubmit={handleSign}>
            <div className="flex items-center border-b border-teal-500 py-2">
                <input
                    className="appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none"
                    type="text" placeholder="to sign" aria-label="Full name" value={hash} onChange={handleChangeHash}/>
                <button
                    className="flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded"
                    type="button" onClick={handleSign}>
                    Sign
                </button>
            </div>
        </form>
             <div>{signature}</div>
        </span>
    )
}
