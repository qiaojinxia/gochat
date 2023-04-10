class TrieNode {
    constructor() {
        this.children = {};
        this.isEndOfWord = false;
    }
}

class Trie {
    constructor() {
        this.root = new TrieNode();
    }

    addWord(word) {
        let current = this.root;
        for (let i = 0; i < word.length; i++) {
            let char = word.charAt(i);
            if (!current.children[char]) {
                current.children[char] = new TrieNode();
            }
            current = current.children[char];
        }
        current.isEndOfWord = true;
    }

    search(word) {
        let current = this.root;
        for (let i = 0; i < word.length; i++) {
            let char = word.charAt(i);
            if (!current.children[char]) {
                return false;
            }
            current = current.children[char];
        }
        return current.isEndOfWord;
    }

    delete(word) {
        let current = this.root;
        this.deleteHelper(current, word, 0);
    }

    deleteHelper(current, word, i) {
        if (i === word.length) {
            current.isEndOfWord = false;
            return Object.keys(current.children).length === 0;
        }
        let char = word.charAt(i);
        if (!current.children[char]) {
            return false;
        }
        let shouldDeleteNode = this.deleteHelper(current.children[char], word, i + 1);
        if (shouldDeleteNode) {
            delete current.children[char];
            return Object.keys(current.children).length === 0;
        }
        return false;
    }

    autoComplete(prefix) {
        let current = this.root;
        for (let i = 0; i < prefix.length; i++) {
            let char = prefix.charAt(i);
            if (!current.children[char]) {
                return [];
            }
            current = current.children[char];
        }
        return this.autoCompleteHelper(current, prefix, []);
    }

    autoCompleteHelper(node, word, result) {
        if (node.isEndOfWord) {
            result.push(word);
        }
        for (let char in node.children) {
            this.autoCompleteHelper(node.children[char], word + char, result);
        }
        return result;
    }
}