import Button from "react-bootstrap/Button";
import InputGroup from "react-bootstrap/InputGroup";
import { Shuffle } from "react-bootstrap-icons";
import AsyncSelect from "react-select/async";
import { MEILI_CLIENT } from "../constants";

// Component for start/end of chain
const StartEnd = ({
  setToSet,
  currentState,
  setState,
  bgVariant,
  placeholder,
}) => {
  return (
    <InputGroup>
      <AsyncSelect
        loadOptions={loadOptions}
        cacheOptions={true}
        defaultOptions={false}
        isClearable={true}
        placeholder={placeholder}
        value={
          currentState === null
            ? null
            : { label: currentState?.name, value: currentState?.id }
        }
        onChange={(n) =>
          setState(n === null ? n : { id: n.value, name: n.label })
        }
        className="flex-fill"
        menuPortalTarget={document.body}
        styles={{
          menuPortal: (base) => ({ ...base, zIndex: 9999 }),
          control: (base, props) => ({
            ...base,
            backgroundColor: `var(--bs-${bgVariant})`,
          }),
          singleValue: (provided) => ({
            ...provided,
            color: "white",
          }),
          placeholder: (provided) => ({
            ...provided,
            color: "white",
          }),
        }}
      />
      {setToSet && (
        <Button variant="secondary" onClick={() => setToSet(() => setState)}>
          <Shuffle />
        </Button>
      )}
    </InputGroup>
  );
};

// load selectable actor options based on search query
export const loadOptions = (inputValue, callback) => {
  setTimeout(() => {
    MEILI_CLIENT.index("actors")
      .search(inputValue, { sort: ["popularity:desc"] })
      .then((resp) => {
        const data = resp.hits.map((h, _) => ({
          value: h?.id,
          label: h?.name,
        }));
        callback(data);
      });
  }, 500);
};

export default StartEnd;
