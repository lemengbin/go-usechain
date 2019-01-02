// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"

	"github.com/usechain/go-usechain/accounts"
	"github.com/usechain/go-usechain/common"
	"github.com/usechain/go-usechain/crypto"
	"github.com/usechain/go-usechain/crypto/ecies"

	"fmt"

	// "github.com/usechain/go-usechain/common/hexutil"
	"github.com/usechain/go-usechain/event"
)

var testSigData = make([]byte, 32)

func TestEncryptDecrypt(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "123456" // not used but required by API
	a1, _, err := ks.NewMainAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	err = ks.Unlock(a1, pass)
	if err != nil {
		t.Error(err)
	}
	pub, err := ks.GetPublicKey(a1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("pubstr: %v", pub))
	pub = pub[2:]

	tmpBytes, err := hex.DecodeString(pub)
	if err != nil {
		t.Error(err)
	}
	publicKey := crypto.ToECDSAPub(tmpBytes)
	pubKey := ecies.ImportECDSAPublic(publicKey)

	tmpBytes, err = ecies.Encrypt(crand.Reader, pubKey, []byte("info"), nil, nil)
	if err != nil {
		t.Error(err)
	}
	infoData := fmt.Sprintf("%x", tmpBytes)
	fmt.Println(fmt.Sprintf("get encrypt result: %v", infoData))

	pri, err := ks.GetPrivateKey(a1)
	if err != nil {
		t.Error(err)
	}
	
	privateKey := ecies.ImportECDSA(pri)
	encryptText, err := hex.DecodeString(infoData)
	if err != nil {
		t.Error(err)
	}
	plaintext, err := privateKey.Decrypt(crand.Reader, encryptText, nil, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("get decrypt result: %v", string(plaintext)))
}
func (ks *KeyStore) getPrivateKey(a accounts.Account, pass string) ([]byte, error) {
	_, key, err := ks.getDecryptedKey(a, pass)
	if err != nil {
		return nil, err
	}
	bytes1 := crypto.FromECDSA(key.PrivateKey)
	return bytes1, nil
}
func TestNewMainAccount(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "123456" // not used but required by API
	a1, abAddr, err := ks.NewMainAccount(pass)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("account.url:  %v\n account.address:  %v", a1.URL, a1.Address.Hex()))
	fmt.Println(fmt.Sprintf("abAddr:%v", abAddr))

	privateBytes, err := ks.getPrivateKey(a1, pass)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(fmt.Sprintf("privateKeyBytes: %v", privateBytes))

	_, key, err := ks.getEncryptedKey(a1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(fmt.Sprintf("use getEncryptedKey to get private key: %v", key))
}

func TestKeyStore(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	auth := "passwd123456"
	a, err := ks.NewAccount(auth)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(a.URL.Path, dir) {
		t.Errorf("account file %s doesn't have dir prefix", a.URL)
	}
	stat, err := os.Stat(a.URL.Path)
	if err != nil {
		t.Fatalf("account file %s doesn't exist (%v)", a.URL, err)
	}
	if runtime.GOOS != "windows" && stat.Mode() != 0600 {
		t.Fatalf("account file has wrong mode: got %o, want %o", stat.Mode(), 0600)
	}
	if !ks.HasAddress(a.Address) {
		t.Errorf("HasAccount(%x) should've returned true", a.Address)
	}
	if err := ks.Unlock(a, auth); err != nil {
		t.Errorf("Unlock error: %v", err)
	}
}

func TestSign(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "" // not used but required by API
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(a1, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.SignHash(accounts.Account{Address: a1.Address}, testSigData); err != nil {
		t.Fatal(err)
	}
}

func TestSignWithPassphrase(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "passwd"
	acc, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	if _, unlocked := ks.unlocked[acc.Address]; unlocked {
		t.Fatal("expected account to be locked")
	}

	_, err = ks.SignHashWithPassphrase(acc, pass, testSigData)
	if err != nil {
		t.Fatal(err)
	}

	if _, unlocked := ks.unlocked[acc.Address]; unlocked {
		t.Fatal("expected account to be locked")
	}

	if _, err = ks.SignHashWithPassphrase(acc, "invalid passwd", testSigData); err == nil {
		t.Fatal("expected SignHashWithPassphrase to fail with invalid password")
	}
}

func TestTimedUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "foo"
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase fails because account is locked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked before unlocking, got ", err)
	}

	// Signing with passphrase works
	if err = ks.TimedUnlock(a1, pass, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// Signing fails again after automatic locking
	time.Sleep(250 * time.Millisecond)
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked timeout expired, got ", err)
	}
}

func TestOverrideUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t, false)
	defer os.RemoveAll(dir)

	pass := "foo"
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Unlock indefinitely.
	if err = ks.TimedUnlock(a1, pass, 5*time.Minute); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// reset unlock to a shorter period, invalidates the previous unlock
	if err = ks.TimedUnlock(a1, pass, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase still works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// Signing fails again after automatic locking
	time.Sleep(250 * time.Millisecond)
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked timeout expired, got ", err)
	}
}

// This test should fail under -race if signing races the expiration goroutine.
func TestSignRace(t *testing.T) {
	dir, ks := tmpKeyStore(t, false)
	defer os.RemoveAll(dir)

	// Create a test account.
	a1, err := ks.NewAccount("")
	if err != nil {
		t.Fatal("could not create the test account", err)
	}

	if err := ks.TimedUnlock(a1, "", 15*time.Millisecond); err != nil {
		t.Fatal("could not unlock the test account", err)
	}
	end := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(end) {
		if _, err := ks.SignHash(accounts.Account{Address: a1.Address}, testSigData); err == ErrLocked {
			return
		} else if err != nil {
			t.Errorf("Sign error: %v", err)
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	t.Errorf("Account did not lock within the timeout")
}

// Tests that the wallet notifier loop starts and stops correctly based on the
// addition and removal of wallet event subscriptions.
func TestWalletNotifierLifecycle(t *testing.T) {
	// Create a temporary kesytore to test with
	dir, ks := tmpKeyStore(t, false)
	defer os.RemoveAll(dir)

	// Ensure that the notification updater is not running yet
	time.Sleep(250 * time.Millisecond)
	ks.mu.RLock()
	updating := ks.updating
	ks.mu.RUnlock()

	if updating {
		t.Errorf("wallet notifier running without subscribers")
	}
	// Subscribe to the wallet feed and ensure the updater boots up
	updates := make(chan accounts.WalletEvent)

	subs := make([]event.Subscription, 2)
	for i := 0; i < len(subs); i++ {
		// Create a new subscription
		subs[i] = ks.Subscribe(updates)

		// Ensure the notifier comes online
		time.Sleep(250 * time.Millisecond)
		ks.mu.RLock()
		updating = ks.updating
		ks.mu.RUnlock()

		if !updating {
			t.Errorf("sub %d: wallet notifier not running after subscription", i)
		}
	}
	// Unsubscribe and ensure the updater terminates eventually
	for i := 0; i < len(subs); i++ {
		// Close an existing subscription
		subs[i].Unsubscribe()

		// Ensure the notifier shuts down at and only at the last close
		for k := 0; k < int(walletRefreshCycle/(250*time.Millisecond))+2; k++ {
			ks.mu.RLock()
			updating = ks.updating
			ks.mu.RUnlock()

			if i < len(subs)-1 && !updating {
				t.Fatalf("sub %d: event notifier stopped prematurely", i)
			}
			if i == len(subs)-1 && !updating {
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	t.Errorf("wallet notifier didn't terminate after unsubscribe")
}

type walletEvent struct {
	accounts.WalletEvent
	a accounts.Account
}

// Tests that wallet notifications and correctly fired when accounts are added
// or deleted from the keystore.
func TestWalletNotifications(t *testing.T) {
	dir, ks := tmpKeyStore(t, false)
	defer os.RemoveAll(dir)

	// Subscribe to the wallet feed and collect events.
	var (
		events  []walletEvent
		updates = make(chan accounts.WalletEvent)
		sub     = ks.Subscribe(updates)
	)
	defer sub.Unsubscribe()
	go func() {
		for {
			select {
			case ev := <-updates:
				events = append(events, walletEvent{ev, ev.Wallet.Accounts()[0]})
			case <-sub.Err():
				close(updates)
				return
			}
		}
	}()

	// Randomly add and remove accounts.
	var (
		live       = make(map[common.Address]accounts.Account)
		wantEvents []walletEvent
	)
	for i := 0; i < 1024; i++ {
		if create := len(live) == 0 || rand.Int()%4 > 0; create {
			// Add a new account and ensure wallet notifications arrives
			account, err := ks.NewAccount("")
			if err != nil {
				t.Fatalf("failed to create test account: %v", err)
			}
			live[account.Address] = account
			wantEvents = append(wantEvents, walletEvent{accounts.WalletEvent{Kind: accounts.WalletArrived}, account})
		} else {
			// Delete a random account.
			var account accounts.Account
			for _, a := range live {
				account = a
				break
			}
			if err := ks.Delete(account, ""); err != nil {
				t.Fatalf("failed to delete test account: %v", err)
			}
			delete(live, account.Address)
			wantEvents = append(wantEvents, walletEvent{accounts.WalletEvent{Kind: accounts.WalletDropped}, account})
		}
	}

	// Shut down the event collector and check events.
	sub.Unsubscribe()
	<-updates
	checkAccounts(t, live, ks.Wallets())
	checkEvents(t, wantEvents, events)
}

// checkAccounts checks that all known live accounts are present in the wallet list.
func checkAccounts(t *testing.T, live map[common.Address]accounts.Account, wallets []accounts.Wallet) {
	if len(live) != len(wallets) {
		t.Errorf("wallet list doesn't match required accounts: have %d, want %d", len(wallets), len(live))
		return
	}
	liveList := make([]accounts.Account, 0, len(live))
	for _, account := range live {
		liveList = append(liveList, account)
	}
	sort.Sort(accountsByURL(liveList))
	for j, wallet := range wallets {
		if accs := wallet.Accounts(); len(accs) != 1 {
			t.Errorf("wallet %d: contains invalid number of accounts: have %d, want 1", j, len(accs))
		} else if accs[0] != liveList[j] {
			t.Errorf("wallet %d: account mismatch: have %v, want %v", j, accs[0], liveList[j])
		}
	}
}

// checkEvents checks that all events in 'want' are present in 'have'. Events may be present multiple times.
func checkEvents(t *testing.T, want []walletEvent, have []walletEvent) {
	for _, wantEv := range want {
		nmatch := 0
		for ; len(have) > 0; nmatch++ {
			if have[0].Kind != wantEv.Kind || have[0].a != wantEv.a {
				break
			}
			have = have[1:]
		}
		if nmatch == 0 {
			t.Fatalf("can't find event with Kind=%v for %x", wantEv.Kind, wantEv.a.Address)
		}
	}
}

func tmpKeyStore(t *testing.T, encrypted bool) (string, *KeyStore) {
	d, err := ioutil.TempDir("", "eth-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	new := NewPlaintextKeyStore
	if encrypted {
		new = func(kd string) *KeyStore { return NewKeyStore(kd, veryLightScryptN, veryLightScryptP) }
	}
	return d, new(d)
}

func TestComputeABKeys(t *testing.T) {
	_, ks := tmpKeyStore(t, true)
	//defer os.RemoveAll(dir)

	// create an account
	auth := "passwd"
	a, err := ks.NewAccount(auth)
	if err != nil {
		t.Fatal(err)
	}

	err = ks.Unlock(a, auth)
	if err != nil {
		t.Errorf("unlock fail. err:%s", err.Error())
	}

	abaccount, abadd, err := ks.NewABaccount(a, auth)
	fmt.Println(abadd)
	if err != nil {
		t.Fatal(err)
	}
	err = ks.Unlock(abaccount, auth)
	if err != nil {
		t.Errorf("unlock fail. err:%s", err.Error())
	}

	fmt.Println("abaccount", abaccount)

	//var abAddr common.ABaddress
	pubbbbbb, err := ks.GetPublicKey(abaccount)
	fmt.Println("pbbbbbbbbbbbbbbbbbbbbbbbbbbb", pubbbbbb)

	abAddr, err := ks.GetABaddr(abaccount)
	if err != nil && len(abAddr) != common.ABaddressLength {
		t.Errorf("Generate waddress error: %v", err)
	}
}

func TestNewABaccount(t *testing.T) {
	_, ks := tmpKeyStore(t, true)

	passwd := "123456"
	acc, err := ks.NewAccount(passwd)
	if err != nil {
		t.Fatal(err)
	}
	err = ks.Unlock(acc, passwd)
	if err != nil {
		t.Error(err)
	}

	ABacc, ABaddr, err := ks.NewABaccount(acc, "123456")

	if err != nil {
		t.Error(err)
	}
	err = ks.Unlock(ABacc, "123456")
	if err != nil {
		t.Error(err)
	}
	abaddr, _ := ks.GetABaddr(ABacc)
	if abaddr != hex.EncodeToString(ABaddr[:]) {
		t.Error("ABaddr is not correct")
	}
}

func TestECDSAPKCompression(t *testing.T) {
	reader := crand.Reader
	curve := btcec.S256()
	ecdsaPrivKey, _ := ecdsa.GenerateKey(curve, reader)
	compressed := ECDSAPKCompression(&ecdsaPrivKey.PublicKey)
	if len(compressed) != 33 {
		t.Error("ECDSAPKCompression is not correct")
	}
}

func TestGenerateBaseABaddress(t *testing.T) {
	reader := crand.Reader
	curve := btcec.S256()
	ecdsaPrivKey, _ := ecdsa.GenerateKey(curve, reader)
	tmp := GenerateBaseABaddress(&ecdsaPrivKey.PublicKey)
	if len(*tmp) != 66 {
		t.Error("GenerateBaseABaddress error")
	}
}
