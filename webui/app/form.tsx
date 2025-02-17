"use client"

import { useState, useEffect } from "react";
import axios from 'axios';
import GitHubButton from 'react-github-btn';
import "./form.css";

interface Option {
  value: string;
  relatedFields?: Field[];
}

interface Field {
  name: string;
  defaultValue: string;
  options: Option[];
  desc: string;
}

interface Wrapper {
  field: Field;
  children: Wrapper[];
}

const ConfigurationForm: React.FC<object> = () => {
  const [formData, setFormData] = useState<Record<string, string>>({});
  const [wrappers, setWrappers] = useState<Wrapper[]>([]);

  useEffect(() => {
    axios.get<{ error: string; fields: Field[] }>('/tts/api/fields')
      .then(response => {
        if (!response || response.status !== 200 || !response.data || response.data.error ) {
          throw 'Error fetching the fields: ' +  response.data.error;
        }
        if (response.data.fields) {
          setWrappers(response.data.fields.map((field) => ({
            field: field,
            children: [],
          })));
        }
      })
      .catch(error => {
        alert('Error fetching the fields: ' +  error);
        console.error('There was an error fetching the fields!', error);
      });
  }, []);

  const handleChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const removeChildren = (targetWrapper: Wrapper, wrappers: Wrapper[]): Wrapper[] => {
    targetWrapper.children.forEach((child) => {
      wrappers = removeChildren(child, wrappers);
      wrappers = wrappers.filter(({ field }) => field !== child.field);
    });
    targetWrapper.children = [];
    return wrappers;
  };

  const handleSelectChange = (targetWrapper: Wrapper, value: string) => {
    const selectedOption = targetWrapper.field.options.find(opt => opt.value === value);
    const relatedFields = selectedOption?.relatedFields;
    setWrappers((prev) => {
      let wrappers = removeChildren(targetWrapper, prev);
      const newChildren = relatedFields ? relatedFields.map((field) => ({ field, children: [] })) : [];
      targetWrapper.children = newChildren;
      wrappers = [...wrappers, ...newChildren];
      return wrappers;
    });
    handleChange(targetWrapper.field.name, value);
  };


  const handleListen = (button: HTMLButtonElement) => {
    // Create a new AudioContext
    const ctx = new window.AudioContext;
  
    // Get the button element to enable/disable it
    button.disabled = true;
  
    // Make the API call to fetch the audio
    axios.post<ArrayBuffer>('/tts/api/invoke', formData, {
        headers: {
          'Content-Type': 'application/json',
        },
        responseType: 'arraybuffer',
      })
      .then(response => {
        if (!response || response.headers["error"] || !response.data ) {
          throw "invoke tts api failed, err:" + response.headers["error"];
        }
        return response.data;
      })
      .then((arrayBuffer: ArrayBuffer) => {
        return ctx.decodeAudioData(arrayBuffer);
      })
      .then((audioBuffer) => {
        // Create a source and connect it to the audio context's destination
        const player = ctx.createBufferSource();
        player.buffer = audioBuffer;
        player.connect(ctx.destination);
        player.start(ctx.currentTime); // Start playback
      })
      .catch(reason => {
        // Handle errors
        alert(`Error: ${reason}`);
      })
      .finally(() => {
        // Re-enable the button after the process is finished
        button.disabled = false;
      });
  };

  const handleGenerateSubscribeURL = () => {
    // Get current host
    const currentHost = window.location.origin;
  
    // Convert formData into URL query parameters
    const newForm = { ...formData };
    delete newForm.text;
    newForm["host"] = currentHost;
    const queryParams = new URLSearchParams(newForm).toString();
  
    // Construct full subscription URL
    const subscribeURL = `${currentHost}/tts/api/subscribe?${queryParams}`;
  
    // Copy URL to clipboard
    navigator.clipboard.writeText(subscribeURL)
      .then(() => {
        alert("Subscription URL copied to clipboard!");
      })
      .catch(err => {
        console.error("Failed to copy URL:", err);
      });
  };
  

  return (
    <div className="config-form">
      <GitHubButton href="https://github.com/chasemao/tts-model-server">Star tts-moder-server</GitHubButton>
      <h2 className="form-title">TTS Model Server WebUI</h2>
      
      <div className="form-group">
        <label className="form-label">Token</label>
        <input className="form-input" onChange={(e) => handleChange("token", e.target.value)} />
        <p className="form-desc">Input -token when running server, can be empty</p>
      </div>
      
      <div className="fields-container">
        {wrappers.map((wrapper) => {
          const { field } = wrapper;
          return (
            <div key={field.name} className="form-group">
              <label className="form-label">{field.name}</label>
              {field.options && field.options.length !== 0 ? (
                <select 
                  className="form-select"
                  value={(() => {
                    if (formData[field.name]) {
                      return formData[field.name];
                    }
                    if (field.defaultValue) {
                      handleSelectChange(wrapper, field.defaultValue);
                      return field.defaultValue
                    }
                    return '';
                  })()}
                  onChange={(value) => handleSelectChange(wrapper, value.target.value)}>
                    {field.options.map(({ value }) => (  
                      <option key={value} value={value}>{value}</option>
                    ))}
                </select>
              ) : (
                <input 
                  className="form-input"
                  defaultValue={field.defaultValue} 
                  onChange={(e) => handleChange(field.name, e.target.value)} 
                />
              )}
            </div>
          );
        })}
      </div>
      
      <div className="form-group">
        <label className="form-label">Test Text</label>
        <input className="form-input" value={(()=>{
          if (formData["text"]) {
            return formData["text"];
          }
          const dv = "It is a paragraph for testing, thank you for staring the github repo.";
          formData["text"] = dv;
          return dv;
        })()}
        onChange={(e) => handleChange("text", e.target.value)} 
        />
        <p className="form-desc">Text for testing</p>
      </div>
      <button className="form-button"
        onClick={(e)=>handleListen(e.currentTarget)}>Listen</button>
      
      <button className="form-button"
        onClick={()=>handleGenerateSubscribeURL()}>Generate subscribe URL</button>
    </div>
  );
};

export default ConfigurationForm;
