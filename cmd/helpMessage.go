// Copyright © 2021 Giannis Fakinos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

var helpMessage = `
 ____                                        ____ _            _    
/ ___|  ___  __ _ _   _  ___ _ __   ___ ___ / ___| | ___   ___| | __
\___ \ / _ \/ _\ | | | |/ _ \ '_ \ / __/ _ \ |   | |/ _ \ / __| |/ /
 ___) |  __/ (_| | |_| |  __/ | | | (_|  __/ |___| | (_) | (__|   < 
|____/ \___|\__, |\__,_|\___|_| |_|\___\___|\____|_|\___/ \___|_|\_\
               |_|                                                  

                     🕔⏳ SequenceClock ⏳🕔
              National Technical University of Athens
            School of Electrical & Computer Engineering            
                 Copyright © 2021 Giannis Fakinos

           
A latency targeting tool for serverless sequences of fuctions. With a
certain time slack provided, SequenceClock will manage to execute a
certain pipeline of actions/functions in the specified time duration. 
Its configuration aims Kubernetes Clusters with Openfaas or Openwhisk
deployed.

Examples of usage:
					
	sc create my-awesome-sequence --list step1, step2, step3
	sc delete my-cool-sequence
`
