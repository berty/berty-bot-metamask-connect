import "../index.css";

const React = require('react');

export default function Display({fieldName, content}: { fieldName: string | null, content: string }) {
    // I don't know css
    return (
        <div className="flex items-center ">
            <div>
                <ul>
                    <li>
                        <strong>{"- "+fieldName}</strong>
                    </li>
                    <li>
                        <span>{content}</span>
                    </li>

                </ul>
            </div>
        </div>
    )
}
