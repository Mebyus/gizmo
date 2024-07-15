function main() {
    let step = document.getElementById("step");
    let panel = new RegistersPanel();

    step.addEventListener("click", () => {
        vmStep().then((state) => {
            panel.changeValues(state.regs);
        });
    });
}

// number of general-purpose registers
const REG_NUM = 4;

class RegistersPanel {
    constructor() {
        this.ip = new Register("ip");
        this.cf = new Register("cf");

        let r = [];
        for (let i = 0; i < REG_NUM; i += 1) {
            r.push(new Register("r" + i.toString()));
        }
        this.r = r;
    }

    changeValues(regs) {
        this.ip.changeValue(regs.ip);
        this.cf.changeValue(regs.cf);

        for (let i = 0; i < REG_NUM; i += 1) {
            this.r[i].changeValue(regs.r[i])
        }
    }
}

/**
 * Send a request to make a single vm step.
 * 
 * @returns {Promise} VM state after the step.
 */
function vmStep() {
    return fetch("/vm/step", {method: "POST"})
        .then((response) => response.json());
}

/**
 * Formats a given byte number as string of hex digits. Values less than 0x10
 * are prefixed with 0, thus function always returns a string with 2 digits.
 * 
 * @param {number} b Byte value. Must be a number in range 0 <= b <= 255. 
 * @returns {string} Byte value in hex format. 
 */
function formatByte(b) {
    if (b < 16) {
        return "0" + b.toString(16);
    }
    return b.toString(16);
}

class Register {
    value;

    // Array of 8 div elements, each representing a single byte.
    // First element represents the most significant byte in register.
    bytes;

    /**
     * @param {string} name Register name.
    */
    constructor(name) {
        let r = document.getElementById(name);
        this.bytes = r.getElementsByClassName("byte");
    }

    /**
     * @param {Array<number>} value Array with 8 numbers. First element
     * corresponds to least significant byte in register.
    */
    changeValue(value) {
        this.value = value;
        
        let j = 8;
        for (let i = 0; i < 8; i += 1) {
            j -= 1;
            this.bytes[j].innerText = formatByte(value[i]);
        }
    }
}

main();
